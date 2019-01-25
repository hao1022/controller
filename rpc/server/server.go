package main

import (
    "../rpc"
    "github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
    for _, config := range rpc.Config {
        for _, route := range config.URLs {
	    url := config.Root + route.URL

	    // GET method
	    if route.Method == "" || route.Method == "get" {
	        router.GET(url, func(c *gin.Context) {
		    context := make(map[string]string)
                    if route.Params != nil {
                        for _, param := range route.Params {
			    context[param] = c.Param(param)
			}
	            }
		    c.Writer.Write(rpc.Get(config, route, context))
		})
	    }

	    // POST method
	    if route.Method == "post" {
	        router.POST(url, func(c *gin.Context) {
		    context := make(map[string]string)
                    if route.Params != nil {
                        for _, param := range route.Params {
			    context[param] = c.Param(param)
			}
		    }
		    c.Writer.Write(rpc.Post(config, route, context))
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
