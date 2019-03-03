package main

import (
    "fmt"
    "strconv"
    "encoding/json"
    "../../../model/tezos"
)

func Transfer(config ConfigType, counter int, amount string, delegate string) {
    //header := tezos.Header()
    //metadata := tezos.Metadata("head")
    count := strconv.Itoa(counter)

    txn := tezos.OperationContentType{
        Kind: "transaction",
	Amount: amount,
	Source: config.Baker,
        Fee: "0",
	Counter: count,
	GasLimit: "11000",
	StorageLimit: "0",
	Destination: delegate}
    txn_str, _ := json.Marshal(txn)
    signature := "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q"

    run_json := fmt.Sprintf("{\"branch\": \"%s\", \"contents\": [%s], \"signature\": \"%s\"}",
                           "BLsiQXDNCTim3SF1RhAPAGQkLT6a9oBbtSSFry49VUte4Q29edW", txn_str, signature)
    fmt.Println(run_json)
    result := tezos.RunOperation(run_json)
    fmt.Println(result)
}
