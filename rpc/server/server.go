package main

import (
    "../rpc"
    "github.com/gin-gonic/gin"
    gorilla_rpc "github.com/gorilla/rpc"
    "github.com/gorilla/rpc/json"
)

var Server *gorilla_rpc.Server

func RegisterService(receiver interface{}, name string) error {
    return Server.RegisterService(receiver, name)
}

func GetServer() *gorilla_rpc.Server {
    if Server != nil {
        return Server
    }

    Server := gorilla_rpc.NewServer()
    Server.RegisterCodec(json.NewCodec(), "application/json")
    Server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

    handler := new(rpc.CommonService)
    Server.RegisterService(handler, "")
    return Server
}

func makeContext(c *gin.Context, params []string, queries []string, ctxt map[string]string) {
    if params != nil {
        for _, param := range params {
            ctxt[param] = c.Param(param)
	}
    }
    if queries != nil {
        for _, query := range queries {
            ctxt[query] = c.Query(query)
	}
    }
}

func RegisterJobs(dispatcher *gin.Engine) {
    for _, config := range rpc.Configurations {
        for _, forward := range config.Forwardings {
	    path := config.Root + forward.Source

	    // GET method
	    if forward.Method == "" || forward.Method == "get" {
		config_, forward_ := config, forward
	        dispatcher.GET(path, func(c *gin.Context) {
                    context := make(map[string]string)
		    makeContext(c, forward_.Params, forward_.Query, context)
		    c.Writer.Write(rpc.Get(config_, forward_, context))
		})
	    }

	    // POST method
	    if forward.Method == "post" {
		config_, forward_ := config, forward
	        dispatcher.POST(path, func(c *gin.Context) {
                    context := make(map[string]string)
		    makeContext(c, forward_.Params, forward_.Query, context)
		    c.Writer.Write(rpc.Post(config_, forward_, context))
		})
	    }

	    // RPC calls
	    if forward.Method == "rpc" {
	        dispatcher.POST(path, func(c *gin.Context) {
		    GetServer().ServeHTTP(c.Writer, c.Request)
		})
            }
        }
    }
}

func main () {
    dispatcher := gin.Default()
    rpc.Initialize()
    RegisterJobs(dispatcher)
    dispatcher.Run()
}
