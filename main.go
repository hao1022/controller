package main

import (
    "os"
    "fmt"
    "strings"
    "bytes"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "github.com/gin-gonic/gin"
)

type URLType struct {
    URL string `json:"url"`
    Source string `json:"source"`
    Method string `json:"method"`
    Data string `json:"data"`
    ContentType string `json:"content_type"`
    Params []string `json:"params"`
}

type ConfigType struct {
    Root string `json:"root"`
    Endpoint string `json:"endpoint"`
    Protocol string `json:"protocol"`
    URLs []URLType `json:"URLs"`
}

func (config *ConfigType) ReadFile(filename string) {
    jsonFile, err := os.Open("config/" + filename)
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal([]byte(byteValue), config)
}

func GetHandler(config ConfigType, route URLType) gin.HandlerFunc {
    return func(c *gin.Context) {
        url := fmt.Sprintf("%s://%s%s", config.Protocol, config.Endpoint, route.Source)
	if route.Params != nil {
	    for _, param := range route.Params {
	        value := c.Param(param)
	        pattern := "$" + param
	        url = strings.Replace(url, pattern, value, -1)
            }
	}
        resp, _ := http.Get(url)
	defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        c.Writer.Write(body)
    }
}

func PostHandler(config ConfigType, route URLType) gin.HandlerFunc {
    return func(c *gin.Context) {
        url := fmt.Sprintf("%s://%s%s", config.Protocol, config.Endpoint, route.Source)
	data_str := route.Data
	if route.Params != nil {
	    for _, param := range route.Params {
	        value := c.Param(param)
	        pattern := "$" + param
	        url = strings.Replace(url, pattern, value, -1)
	        data_str = strings.Replace(data_str, pattern, value, -1)
	    }
	}
	fmt.Println(data_str)
	data := []byte(data_str)
	resp, _ := http.Post(url, route.ContentType, bytes.NewBuffer(data))
	defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        c.Writer.Write(body)
    }
}

func RegisterConfig(router *gin.Engine, config_file string) {
    var config ConfigType

    config.ReadFile(config_file)

    for _, route := range config.URLs {
	url := config.Root + route.URL
	if route.Method == "" || route.Method == "get" {
	    router.GET(url, GetHandler(config, route))
	}
	if route.Method == "post" {
	    router.POST(url, PostHandler(config, route))
	}
    }
}

func main () {
    router := gin.Default()
    RegisterConfig(router, "gxchain.json")
    RegisterConfig(router, "eos.json")
    RegisterConfig(router, "tezos.json")
    router.Run()
}
