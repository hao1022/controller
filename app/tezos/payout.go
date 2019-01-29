package main

import (
    "fmt"
    "../../model/tezos"
    "encoding/json"
)

func main() {
    tezos.Initialize()
    c := tezos.CurrentLevel()
    cc, _ := json.Marshal(c)
    fmt.Println(string(cc))
}
