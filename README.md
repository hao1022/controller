## Install dependencies

```sh
$ go get -u github.com/gin-gonic/gin
$ go get -u github.com/gorilla/rpc
```

## Start server

```sh
$ go run main.go
```

## Query

```sh
$ curl -X POST localhost:8080/gxchain/get_full_accounts
$ curl localhost:8080/eos/get_info
$ curl localhost:8080/tezos/balance/head/tz1awXW7wuXy21c66vBudMXQVAPgRnqqwgTH
```
