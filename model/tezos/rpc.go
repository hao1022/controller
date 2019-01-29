package tezos

import (
    _ "fmt"
    "../../rpc/rpc"
    "encoding/json"
)

func Initialize() {
    rpc.Initialize()
}

func CurrentLevel() CurrentLevelType {
    var current_level CurrentLevelType
    body := rpc.Get(rpc.Config["tezos"], *rpc.Config["tezos"].Indices["/current_level"], nil)
    json.Unmarshal(body, &current_level)
    return current_level
}

func Header() BlockHeaderType {
    var blockheader BlockHeaderType
    body := rpc.Get(rpc.Config["tezos"], *rpc.Config["tezos"].Indices["/blockheader"], nil)
    json.Unmarshal(body, &blockheader)
    return blockheader
}


