package main

import (
    "fmt"
    "../../../model/tezos"
    "strconv"
)

type BakerRewardType struct {
    SelfReward int
    FeeReward int
    LooseReward int
    TotalReward int
}

type RewardType struct {
    BakerRewards     BakerRewardType
    Delegators       []string
    DelegatorRewards []int
    StakingBalance   int
}

type StolenBlockType struct {
    Level    int
    Hash     string
    Priority int
    Reward   int
    Fees     int
}

func SnapshotHeight(cycle int, snapshot int, cycle_length int,
                    snapshot_interval int) int {
    return (cycle - 7) * cycle_length + ((snapshot + 1) * snapshot_interval)
}

func BlockHashByLevel(level int) string {
    var past_head string

    head_header := tezos.Header()
    past_head = fmt.Sprintf("%s~%d", head_header.Hash, head_header.Level - level)
    past_header := tezos.HeaderAt(past_head)
    if (head_header.Level != past_header.Level) {
        fmt.Println("should not happen: tezos rpc fault, wrong level")
    }
    return past_header.Hash
}

func HashToQuery(cycle int, cycle_length int) string {
    var level_to_query int

    header := tezos.Header()
    current_level := tezos.CurrentLevelAt(header.Hash)
    blocks_ago := cycle_length * (current_level.Cycle - cycle)
    if header.Level - blocks_ago < header.Level {
        level_to_query = header.Level - blocks_ago
    } else {
	level_to_query = header.Level
    }
    return BlockHashByLevel(level_to_query)
}

func SnapshotHash(cycle int, cycle_length int, snapshot_interval int) string {
    hash := HashToQuery(cycle, cycle_length)
    cycle_info := tezos.CycleInfo(hash, cycle)
    block_height := SnapshotHeight(cycle, cycle_info.Snapshot, cycle_length, snapshot_interval)
    return BlockHashByLevel(block_height)
}

func SnapshotLevel(cycle int, cycle_length int, snapshot_interval int) int {
    hash := HashToQuery(cycle, cycle_length)
    cycle_info := tezos.CycleInfo(hash, cycle)
    return SnapshotHeight(cycle, cycle_info.Snapshot, cycle_length, snapshot_interval)
}

func GetContributingBalancesFor(cycle_length int, snapshot_interval int, cycle int, delegate string) (int, []string, []int) {
    var total_balance int
    var total_frozen_reward int
    var balances []int

    snapshot_block_hash := SnapshotHash(cycle, cycle_length, snapshot_interval)

    delegators := tezos.DelegatedContracts(snapshot_block_hash, delegate)
    for _, delegator := range delegators {
        balance_str := tezos.BalanceAt(snapshot_block_hash, delegator)
	balance, _ := strconv.Atoi(balance_str)
	total_balance += balance
	balances = append(balances, balance)
    }

    full_balance_str := tezos.DelegateBalanceAt(snapshot_block_hash, delegate)
    full_balance, _ := strconv.Atoi(full_balance_str)

    frozen_balance := tezos.FrozenBalanceByCycle(snapshot_block_hash, delegate)
    for _, frozen_reward := range frozen_balance {
	reward, _ := strconv.Atoi(frozen_reward.Rewards)
	total_frozen_reward += reward
    }

    self_balance := full_balance - total_frozen_reward

    staking_balance_str := tezos.StakingBalanceAt(snapshot_block_hash, delegate)
    staking_balance, _ := strconv.Atoi(staking_balance_str)
    if self_balance + total_balance != staking_balance {
        fmt.Println("should not happen")
    }

    return staking_balance, delegators, balances
}

func EstimatedRewards(cycle_length int, cycle int, delegate string) int{
    hash := HashToQuery(cycle, cycle_length)
    baking_rights := tezos.BakingRightsFor(hash, delegate, cycle)
    endorsing_rights := tezos.EndorsingRightsFor(hash, delegate, cycle)

    num_baking_rights := 0
    for _, bright := range baking_rights {
        if bright.Priority != 0 {
            num_baking_rights += 1
        }
    }

    num_endorsing_slots := 0
    for _, eright := range endorsing_rights {
        num_endorsing_slots += len(eright.Slots)
    }

    total_rewards := 16 * num_baking_rights + 2 * num_endorsing_slots
    return total_rewards
}

func CalculateRewardsFor(cycle_length int, snapshot_interval int, cycle int, delegate string, rewards int, fee_percent int) RewardType {
    staking_balance, delegators, balances := GetContributingBalancesFor(cycle_length,
                                             snapshot_interval, cycle, delegate)
    total_balance := 0
    for _, balance := range balances[1:] {
        total_balance += balance
    }
    baker_balance := balances[0]
    baker_self_reward := baker_balance * rewards / total_balance
    baker_fee_reward := fee_percent * rewards * (total_balance - baker_balance) / ( total_balance * 100 )
    var delegator_rewards []int
    var total_delegator_rewards int
    for _, balance := range balances[1:] {
	reward := balance * (100 - fee_percent) * rewards / ( total_balance * 100 )
        delegator_rewards = append(delegator_rewards, reward)
	total_delegator_rewards += reward
    }
    baker_loose_reward := rewards - total_delegator_rewards - baker_self_reward - baker_fee_reward
    baker_total_reward := baker_self_reward + baker_fee_reward + baker_loose_reward
    if baker_total_reward + total_delegator_rewards != rewards {
        fmt.Println("should not happen: rewards mismatch")
    }

    baker_rewards := BakerRewardType{baker_self_reward,
                                   baker_fee_reward,
				   baker_loose_reward,
				   baker_total_reward}

    return RewardType{baker_rewards, delegators, delegator_rewards, staking_balance}
}

func StolenBlocks(cycle_length int, cycle int, delegate string) []StolenBlockType {
    hash := HashToQuery(cycle_length, cycle)
    baking_rights := tezos.BakingRightsFor(hash, delegate, cycle)
    var stolen_blocks []StolenBlockType

    for _, right := range baking_rights {
        if right.Priority >= 0 {
            priority := right.Priority
	    level := right.Level

            hash := BlockHashByLevel(level)
            metadata := tezos.Metadata(hash)
	    if metadata.Baker != delegate {
                continue
            } else {
                operations := tezos.Operations(hash)
		fees := 0
		var reward int
		for _, content := range operations.Contents {
		    if content.Fee != "" {
                        fee, _ := strconv.Atoi(content.Fee)
			fees += fee
		    }
		}
		for _, update := range metadata.BalanceUpdates {
                    if update.Kind == "freezer" && update.Category == "rewards" && update.Delegate == delegate {
                        reward, _ = strconv.Atoi(update.Change)
		    }
                }
                stolen_blocks = append(stolen_blocks, StolenBlockType{level, hash, priority, reward, fees})
	    }
	}
    }
    return stolen_blocks
}

func main() {
    tezos.Initialize()
}
