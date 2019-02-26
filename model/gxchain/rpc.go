package gxchain

import (
    "strconv"
    "../../rpc/rpc"
    "encoding/json"
)

func Initialize() {
    rpc.Initialize()
}

func WitnessByAccount(account string) WitnessType {
    var witness WitnessType
    body := rpc.Post(rpc.Config["gxchain"],
                    *rpc.Config["gxchain"].Indices["/get_witness_by_account/:account"],
                    map[string]string{"account": account})
    json.Unmarshal(body, &witness)
    return witness
}
