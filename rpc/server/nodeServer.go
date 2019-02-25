package main

import (
    "../node"
    "net/http"
    gorilla_rpc "github.com/gorilla/rpc"
    "github.com/gorilla/rpc/json"
)

var Server *gorilla_rpc.Server


func initServer() *gorilla_rpc.Server {
    if Server != nil {
        return Server
    }

    Server := gorilla_rpc.NewServer()
    Server.RegisterCodec(json.NewCodec(), "application/json")
    Server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
    Server.RegisterService(new(node.Action), "")
    return Server
}


func main () {
  server := initServer()
  http.Handle("/rpc", server)
  http.ListenAndServe(":9090", nil)
}
