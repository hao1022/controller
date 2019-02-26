package main

import (
    "../controller"
    "github.com/gin-gonic/gin"
)

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
    for _, config := range controller.Configurations {
        for _, forward := range config.Forwardings {
	    path := config.Root + forward.Source

	    // GET method
	    if forward.Method == "" || forward.Method == "get" {
		config_, forward_ := config, forward
	        dispatcher.GET(path, func(c *gin.Context) {
                    context := make(map[string]string)
		    makeContext(c, forward_.Params, forward_.Query, context)
		    c.Writer.Write(controller.Get(config_, forward_, context))
		})
	    }

	    // POST method
	    if forward.Method == "post" {
		config_, forward_ := config, forward
	        dispatcher.POST(path, func(c *gin.Context) {
                    context := make(map[string]string)
		    makeContext(c, forward_.Params, forward_.Query, context)
		    c.Writer.Write(controller.Post(config_, forward_, context))
		})
	    }

	    // Action calls such as backup and upgrade, not test yet
	    if forward.Method == "rpc" {
		config_, forward_ := config, forward
	        dispatcher.POST(path, func(c *gin.Context) {
                    context := make(map[string]string)
		    makeContext(c, forward_.Params, forward_.Query, context)
		    c.Writer.Write(controller.Post(config_, forward_, context))
		})
            }
        }
    }
}

func main () {
    dispatcher := gin.Default()
    controller.Initialize("../../config/forward/")
    RegisterJobs(dispatcher)
    dispatcher.Run()
}
