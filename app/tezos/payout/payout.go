package main

import (
    "fmt"
    "../../../model/tezos"
    "strconv"
    "os"
    "io"
    "os/exec"
    "strings"
    "bufio"
)

type ConfigType struct {
    TezosClientPath   string
    Baker             string
    CycleLength       int
    SnapshotInterval  int
    Delegate          string
    DelegateName      string
    FeePercent        int
    StartingCycle     int
}

type StolenBlockType struct {
    Level    int
    Hash     string
    Priority int
    Reward   int
    Fees     int
}

type CycleRewardType struct {
    Realized            int
    Paid                int
    RealizedDifference  int
    EstimatedDifference int
}

type BakerRewardType struct {
    SelfReward int
    FeeReward int
    TotalReward int
}

type DelegatorPayoutType struct {
  Balance             int
  EstimatedRewards    int
  FinalRewards        int
  PayoutOperationHash string
}

func GetEstimatesForCycle(config ConfigType, cycle int) RewardType {

    cycle_length := config.CycleLength
    snapshot_interval := config.SnapshotInterval
    baker := config.Baker
    fee_percent := config.FeePercent

    estimated_rewards := EstimatedRewards(cycle_length, snapshot_interval, baker)
    rewards := CalculateRewardsFor(cycle_length, snapshot_interval, cycle, baker,
                                   estimated_rewards, fee_percent)
    return rewards
}

func GetEstimates(config ConfigType) []RewardType {
    var rewards []RewardType
    current_level := tezos.CurrentLevel()
    current_cycle := current_level.Cycle
    known_cycle := current_cycle + 5
    for cycle := current_cycle; cycle <= known_cycle; cycle++ {
        rewards = append(rewards, GetEstimatesForCycle(config, cycle))
    }
    return rewards
}

func GetActualsForCycle(config ConfigType, cycle int) RewardType {
    cycle_length := config.CycleLength
    snapshot_interval := config.SnapshotInterval
    baker := config.Baker
    fee_percent := config.FeePercent

    var reward RewardType

    //stolen := StolenBlocks(cycle_length, cycle, baker)
    hash := HashToQuery(cycle + 1, cycle_length)
    frozen_balance_by_cycle := tezos.FrozenBalanceByCycle(hash, baker)
    for _, balance := range frozen_balance_by_cycle {
        if balance.Cycle == cycle {
	    fee_rewards, _ := strconv.Atoi(balance.Fees)
	    //extra_rewards := fee_rewards
	    balance_rewards, _ := strconv.Atoi(balance.Rewards)
	    realized_rewards := fee_rewards + balance_rewards
	    //estimated_rewards := EstimatedRewards(cycle_length, cycle, baker)
	    //paid_rewards := estimated_rewards + extra_rewards
	    //realized_difference := realized_rewards - paid_rewards
	    //estimated_difference := estimated_rewards - paid_rewards

            reward = CalculateRewardsFor(cycle_length, snapshot_interval, cycle, baker, realized_rewards, fee_percent)
	}
    }
    return reward
}

func GetActuals(config ConfigType) []RewardType {
    var rewards []RewardType
    current_level := tezos.CurrentLevel()
    current_cycle := current_level.Cycle
    known_cycle := current_cycle + 5
    for cycle := current_cycle-6; cycle <= known_cycle; cycle++ {
        rewards = append(rewards, GetActualsForCycle(config, cycle))
    }
    return rewards
}

var Config ConfigType = ConfigType{
    "/home/ubuntu/tezos/tezos-client", // path to tezos-client
    "tz1awXW7wuXy21c66vBudMXQVAPgRnqqwgTH", // baker account
    4096, // cycle length, const
    256,  // snapshot interval, const
    "tz1awXW7wuXy21c66vBudMXQVAPgRnqqwgTH", // delegate account
    "infstones", // delegate name
    10, // fee percent, 10% by default
    64} // starting cycle

func PrintRewards(rewards []RewardType) {
    for _, reward := range rewards {
        fmt.Printf("Cycle: %d\n", reward.Cycle)
        fmt.Printf("Baker Self Reward: %d\n", reward.BakerRewards.SelfReward)
        fmt.Printf("Baker Fee Reward: %d\n", reward.BakerRewards.FeeReward)
        fmt.Printf("Baker Total Reward: %d\n", reward.BakerRewards.TotalReward)
        for i, _ := range reward.Delegators {
                fmt.Println(reward.Delegators[i])
                fmt.Printf("  Balance: %d\n", reward.DelegatorBalances[i])
                fmt.Printf("  RawReward: %d\n", reward.DelegatorRawRewards[i])
                fmt.Printf("  Reward: %d\n", reward.DelegatorRewards[i])
                fmt.Printf("  Share: %.2f%%\n", reward.DelegatorShares[i])
        }
        fmt.Printf("Staking Balance: %d\n", reward.StakingBalance)
        fmt.Printf("Total Reward: %d\n", reward.TotalReward)
    }
}

func Payout(rewards []RewardType) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("Password:")
    password, _ := reader.ReadString('\n')

    for _, reward := range rewards {
	fmt.Printf("Payout cycle %d? [y/n]:", reward.Cycle)
        text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	fmt.Println(text)
        if text != "y" {
	    continue
	}
        for i, _ := range reward.Delegators {

		// format command
		if reward.DelegatorRewards[i] == 0 {
		    continue
		}
		amount := float64(reward.DelegatorRewards[i]) / 1000000.0
		amount_str := strconv.FormatFloat(amount, 'g', 6, 64)
		cmd := fmt.Sprintf("%s transfer %s from %s to %s", Config.TezosClientPath,
		                   amount_str, Config.Delegate, reward.Delegators[i])
		// print out command
                fmt.Println(cmd)

		// execute command
		process := exec.Command(Config.TezosClientPath,
		                   "transfer", amount_str,
				   "from", Config.Delegate, "to", reward.Delegators[i])
                stdin, err := process.StdinPipe()
		if err != nil {
                    fmt.Println(err)
                }
                defer stdin.Close()
		process.Stdout = os.Stdout
		process.Stderr = os.Stderr
		if err = process.Start(); err != nil {
                    fmt.Println("An error occured: ", err)
		}
                io.WriteString(stdin, password)
		process.Wait()
        }
    }
}

func main() {
    tezos.Initialize()
    rewards := GetActuals(Config)
    Payout(rewards)
}
