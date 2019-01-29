package tezos

import (
    _ "fmt"
    "../../rpc/rpc"
    "encoding/json"
)

func Initialize() {
    rpc.Initialize()
}

func CurrentLevelAt(hash string) CurrentLevelType {
    var current_level CurrentLevelType
    body := rpc.Get(rpc.Config["tezos"], *rpc.Config["tezos"].Indices["/current_level"],
                    map[string]string{"head": hash})
    json.Unmarshal(body, &current_level)
    return current_level
}

func CurrentLevel() CurrentLevelType {
    return CurrentLevelAt("head")
}

func HeaderAt(hash string) BlockHeaderType {
    var blockheader BlockHeaderType
    body := rpc.Get(rpc.Config["tezos"], *rpc.Config["tezos"].Indices["/blockheader"],
                    map[string]string{"head": hash})
    json.Unmarshal(body, &blockheader)
    return blockheader
}

func Header() BlockHeaderType {
    return HeaderAt("head")
}

