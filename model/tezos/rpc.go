package tezos

import (
    "strconv"
    "../../rpc/rpc"
    "encoding/json"
)

func Initialize() {
    rpc.Initialize()
}

func CurrentLevelAt(hash string) CurrentLevelType {
    var current_level CurrentLevelType
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/current_level/:head"],
                    map[string]string{"head": hash})
    json.Unmarshal(body, &current_level)
    return current_level
}

func CurrentLevel() CurrentLevelType {
    return CurrentLevelAt("head")
}

func HeaderAt(hash string) BlockHeaderType {
    var blockheader BlockHeaderType
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/blockheader/:head"],
                    map[string]string{"head": hash})
    json.Unmarshal(body, &blockheader)
    return blockheader
}

func Header() BlockHeaderType {
    return HeaderAt("head")
}

func CycleInfo(hash string, cycle int) CycleInfoType {
    var cycle_info CycleInfoType
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/cycle_info/:head/:cycle"],
                    map[string]string{
			    "head": hash,
			    "cycle": strconv.Itoa(cycle)})
    json.Unmarshal(body, &cycle_info)
    return cycle_info
}

func DelegatedContracts(hash string, delegate string) []string {
    var contracts []string
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/delegated_contracts/:head/:account"],
                    map[string]string{
			    "head": hash,
			    "account": delegate})
    json.Unmarshal(body, &contracts)
    return contracts
}

func BalanceAt(hash string, delegate string) string {
    var balance string
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/balance/:head/:account"],
                    map[string]string{
			    "head": hash,
			    "account": delegate})
    json.Unmarshal(body, &balance)
    return balance
}

func DelegateBalanceAt(hash string, delegate string) string {
    var balance string
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/delegate_balance/:head/:account"],
                    map[string]string{
			    "head": hash,
			    "account": delegate})
    json.Unmarshal(body, &balance)
    return balance
}

func FrozenBalanceByCycle(hash string, delegate string) []FrozenBalanceByCycleType {
    var balance []FrozenBalanceByCycleType
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/frozen_balance_by_cycle/:head/:account"],
                    map[string]string{
			    "head": hash,
			    "account": delegate})
    json.Unmarshal(body, &balance)
    return balance
}

func StakingBalanceAt(hash string, delegate string) string {
    var balance string
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/staking_balance/:head/:account"],
                    map[string]string{
			    "head": hash,
			    "account": delegate})
    json.Unmarshal(body, &balance)
    return balance
}

func BakingRightsFor(hash string, delegate string, cycle int) []BakingRightType {
    var baking_rights []BakingRightType
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/baking_rights/:head/:account/:cycle"],
                    map[string]string{
			    "head": hash,
			    "account": delegate,
		            "cycle": strconv.Itoa(cycle)})
    json.Unmarshal(body, &baking_rights)
    return baking_rights
}

func EndorsingRightsFor(hash string, delegate string, cycle int) []EndorsingRightType {
    var endorsing_rights []EndorsingRightType
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/endorsing_rights/:head/:account/:cycle"],
                    map[string]string{
			    "head": hash,
			    "account": delegate,
		            "cycle": strconv.Itoa(cycle)})
    json.Unmarshal(body, &endorsing_rights)
    return endorsing_rights
}

func Metadata(hash string) BlockMetadataType {
    var metadata BlockMetadataType
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/metadata/:head"],
                    map[string]string{"head": hash})
    json.Unmarshal(body, &metadata)
    return metadata
}

func Operations(hash string) OperationType {
    var operations OperationType
    body := rpc.Get(rpc.Config["tezos"],
                    *rpc.Config["tezos"].Indices["/operations/:head"],
                    map[string]string{"head": hash})
    json.Unmarshal(body, &operations)
    return operations
}
