package main

import (
    "fmt"
    "strconv"
    "encoding/json"
    "../../../model/tezos"

//    "io"
//    "os"
//    "bytes"
    "strings"
    "bufio"
    "time"
    "os/exec"
    "github.com/kr/pty"
    "github.com/btcsuite/btcutil/base58"
)

func Sign(config ConfigType, account string, operation string) string {
    signature := ""
    operation = "55a3af1a962a22083d632b97fe2063f09b8c393131fb8f55da59b1c811eb4eed080000a7d901fea5daf54859d1692593d72b558f34d64900d39f07f85500e50201f0c38bb6fd4c6a2ec16f1f7f306d9c8ae6b68cda0000"
    op_bytes := fmt.Sprintf("0x03%s", operation)
    process := exec.Command(
	    config.TezosClientPath, "-A",
	    config.Endpoint, "sign", "bytes", op_bytes,
	    "for", account)

    //fmt.Println(config.TezosClientPath)
    //fmt.Println(config.Endpoint)
    //fmt.Println(op_bytes)
    //fmt.Println(account)
    tty, _ := pty.Start(process)
    defer tty.Close()

    time.Sleep(1000 * time.Millisecond)

    // redirect tty stdin
    go func() {
        //tty.Write([]byte(config.Password + "\n"))
        tty.Write([]byte("hao4infinitystones\n"))
    }()
    time.Sleep(100 * time.Millisecond)
    scanner := bufio.NewScanner(tty)
    for scanner.Scan() {
        //fmt.Println(scanner.Text())
	line := scanner.Text()
	if strings.HasPrefix(line, "Signature: ") {
	    signature = line[11:]
	    return signature
	}
    }
    return signature
}

func Transfer(config ConfigType, counter int, amount string, delegate string) {
    header := tezos.Header()
    //metadata := tezos.Metadata("head")
    count := strconv.Itoa(counter)

    txn := tezos.OperationContentType{
        Kind: "transaction",
	Amount: amount,
//	Source: config.Baker,
	Source: "tz1SEj1r3cNz2v1NRCUTKPvcd9wRxANytuNv",
        Fee: "0",
	Counter: count,
	GasLimit: "11000",
	StorageLimit: "0",
	Destination: delegate}
    txn_str, _ := json.Marshal(txn)
    //signature := "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q"

    //run_json := fmt.Sprintf("{\"branch\": \"%s\", \"contents\": [%s], \"signature\": \"%s\"}",
    //                       header.Hash, txn_str, signature)
    //fmt.Println(run_json)
    //run_json_result := tezos.RunOperation(run_json)
    //if run_json_result.Contents == nil || run_json_result.Contents[0].Metadata.Result.Status != "applied" {
    //    fmt.Println(run_json_result)
    //    return
    //}

    sign_json := fmt.Sprintf("{\"branch\": \"%s\", \"contents\": [%s]}", header.Hash, txn_str)
    //fmt.Println(sign_json)
    sign_json_result := tezos.ForgeOperations(sign_json)
    //fmt.Println(sign_json_result)

    sig := Sign(config, "hao1022", sign_json_result)
    fmt.Println(sig)

    base16sig := base58.Decode(sig)
    truesig := base16sig[5:len(base16sig)-4]
    fmt.Printf("%x\n", truesig)

    signed_op := fmt.Sprintf("%s%s", sign_json_result, truesig)

    preapply_json := fmt.Sprintf(
	    "[{\"protocol\": \"%s\", \"branch\": \"%s\", \"contents\": [%s], \"signature\": \"%s\"}]",
	    header.Protocol, header.Hash, txn_str, sig)
    fmt.Println(preapply_json)
    preapply_json_result := tezos.PreapplyOperations(preapply_json)
    fmt.Println(preapply_json_result)

    txn_hash := tezos.Injection(signed_op)
    fmt.Println(txn_hash)
}
