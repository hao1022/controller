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

func RegisterRoutes(router *gin.Engine) {
    for _, config := range rpc.Config {
        for _, route := range config.URLs {
	    url := config.Root + route.URL

	    // GET method
	    if route.Method == "" || route.Method == "get" {
		config_, route_ := config, route
	        router.GET(url, func(c *gin.Context) {
                    context := make(map[string]string)
		    makeContext(c, route_.Params, route_.Query, context)
		    c.Writer.Write(rpc.Get(config_, route_, context))
		})
	    }

	    // POST method
	    if route.Method == "post" {
		config_, route_ := config, route
	        router.POST(url, func(c *gin.Context) {
                    context := make(map[string]string)
		    makeContext(c, route_.Params, route_.Query, context)
		    c.Writer.Write(rpc.Post(config_, route_, context))
		})
	    }

	    // RPC calls
	    if route.Method == "rpc" {
	        router.POST(url, func(c *gin.Context) {
		    GetServer().ServeHTTP(c.Writer, c.Request)
		})
            }
        }
    }
}

func main () {
    router := gin.Default()
    rpc.Initialize()
    RegisterRoutes(router)
    router.Run()
}
