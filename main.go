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
    Source string `json:"source"`
    Target string `json:"target"`
    Method string `json:"method"`
    Data string `json:"data"`
    ConType string `json:"contype"`
    Params []string `json:"params"`
}

type ConfigType struct {
    Root string `json:"root"`
    Endpoint string `json:"endpoint"`
    Protocol string `json:"protocol"`
    Urls []URLType `json:"urls"`
}

func (config *ConfigType) ReadConfig(filename string) {
    str := []string{"config/", filename}
    fullpath_filename := strings.Join(str, "")

    jsonFile, err := os.Open(fullpath_filename)
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal([]byte(byteValue), config)
}

func GetHandler(url URLType, target string) gin.HandlerFunc {
    return func(c *gin.Context) {
	for _, param := range url.Params {
	    value := c.Param(param)
	    pattern := "$" + param
	    target = strings.Replace(target, pattern, value, -1)
        }
        resp, _ := http.Get(target)
	defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        c.Writer.Write(body)
    }
}

func PostHandler(url URLType, target string, data_str string) gin.HandlerFunc {
    return func(c *gin.Context) {
	for _, param := range url.Params {
	    value := c.Param(param)
	    pattern := "$" + param
	    target = strings.Replace(target, pattern, value, -1)
	    data_str = strings.Replace(data_str, pattern, value, -1)
	}
	data := []byte(data_str)
	resp, _ := http.Post(target, url.ConType, bytes.NewBuffer(data))
	defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        c.Writer.Write(body)
    }
}

func Bind(router *gin.Engine, config_file string) {
    var config ConfigType

    config.ReadConfig(config_file)
    fmt.Println(config)

    for _, e := range config.Urls {
	source := config.Root + e.Source
	target := fmt.Sprintf("%s://%s%s", config.Protocol, config.Endpoint, e.Target)
	if e.Method == "get" {
	    router.GET(source, GetHandler(e, target))
	}
	if e.Method == "post" {
	    router.POST(source, PostHandler(e, target, e.Data))
	}
    }
}

func main () {
    router := gin.Default()
    Bind(router, "gxchain.json")
    Bind(router, "eos.json")
    Bind(router, "tezos.json")
    router.Run()
}
