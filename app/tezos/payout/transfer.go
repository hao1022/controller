package main

import (
    "fmt"
    "strconv"
    "encoding/json"
    "../../../model/tezos"

    "strings"
    "bufio"
    "time"
    "os/exec"
    "github.com/kr/pty"
    "github.com/btcsuite/btcutil/base58"
)

func Sign(config ConfigType, operation string) string {
    signature := ""
    op_bytes := fmt.Sprintf("0x03%s", operation)
    process := exec.Command(
	    config.TezosClientPath, "-A",
	    config.Endpoint, "sign", "bytes", op_bytes,
	    "for", config.Account)

    tty, _ := pty.Start(process)
    defer tty.Close()

    time.Sleep(1000 * time.Millisecond)

    // redirect tty stdin
    go func() {
        tty.Write([]byte(config.Password + "\n"))
    }()
    time.Sleep(100 * time.Millisecond)
    scanner := bufio.NewScanner(tty)
    for scanner.Scan() {
	line := scanner.Text()
	if strings.HasPrefix(line, "Signature: ") {
	    signature = line[11:]
	    return signature
	}
    }
    return signature
}

func Transfer(config ConfigType, counter int, amount string, delegate string) string {
    header := tezos.Header()
    count := strconv.Itoa(counter)

    txn := tezos.OperationContentType{
        Kind: "transaction",
	Amount: amount,
	Source: config.Baker,
        Fee: "2300",
	Counter: count,
	GasLimit: "11000",
	StorageLimit: "0",
	Destination: delegate}
    txn_str, _ := json.Marshal(txn)
    signature := "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q"

    run_json := fmt.Sprintf("{\"branch\": \"%s\", \"contents\": [%s], \"signature\": \"%s\"}",
                           header.Hash, txn_str, signature)
    run_json_result := tezos.RunOperation(run_json)
    if run_json_result.Contents == nil || run_json_result.Contents[0].Metadata.Result.Status != "applied" {
	fmt.Println("RunOperation failed:")
        fmt.Println(run_json)
        return ""
    }

    sign_json := fmt.Sprintf("{\"branch\": \"%s\", \"contents\": [%s]}", header.Hash, txn_str)
    sign_json_result := tezos.ForgeOperations(sign_json)

    sig := Sign(config, sign_json_result)
    if sig == "" {
	fmt.Println("Signing failed")
        fmt.Println(run_json)
	return ""
    }

    base16sig := base58.Decode(sig)
    truesig := base16sig[5:len(base16sig)-4]

    signed_op := fmt.Sprintf("\"%s%x\"", sign_json_result, truesig)

    preapply_json := fmt.Sprintf(
	    "[{\"protocol\": \"%s\", \"branch\": \"%s\", \"contents\": [%s], \"signature\": \"%s\"}]",
	    header.Protocol, header.Hash, txn_str, sig)
    preapply_json_result := tezos.PreapplyOperations(preapply_json)
    if preapply_json_result == nil {
	fmt.Println("Preapply failed:")
	fmt.Println(run_json)
	return ""
    }

    txn_hash := tezos.Injection(signed_op)
    return txn_hash
}
