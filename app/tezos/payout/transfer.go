package main

import (
    "fmt"
    "encoding/json"
    "../../../model/tezos"
)

func Transfer(config ConfigType, amount string, delegate string, reward string) {
    counter := tezos.Counter(config.Delegate)
    header := tezos.Header()
    //metadata := tezos.Metadata("head")

    txn := tezos.OperationContentType{
        Kind: "transaction",
	Amount: reward,
	Source: config.Baker,
        Fee: "0",
	Counter: counter,
	GasLimit: "11000",
	StorageLimit: "0",
	Destination: Config.Delegate}
    txn_str, _ := json.Marshal(txn)
    signature := "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q"

    run_json := fmt.Sprintf("{\"branch\": \"%s\", \"contents\": %s, \"signature\": \"%s\"}",
                           header.Hash, txn_str, signature)
    fmt.Println(run_json)
    result := tezos.RunOperation(run_json)
    fmt.Println(result)
}
