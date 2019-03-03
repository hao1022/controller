package controller

import (
    "os"
    "fmt"
    "bytes"
    "strings"
    "io/ioutil"
    "net/http"
    "encoding/json"
)

type Forwarding struct {
    Source string `json:"source"`
    Target string `json:"target"`
    Host string `json:"host"`
    Method string `json:"method"`
    Data string `json:"data"`
    ContentType string `json:"content_type"`
    Params []string `json:"params"`
    Query []string `json:"query"`
}

type Config struct {
    Root string `json:"root"`
    Host string `json:"host"`
    Protocol string `json:"protocol"`
    Forwardings []Forwarding `json:"forwardings"`
    Indices map[string]*Forwarding
}

var Configurations map[string]Config


func Get(config Config, forward Forwarding, context map[string]string) []byte {
    target := forward.Target
    if context != nil {
        for param, value := range context {
            pattern := "$" + param
            target = strings.Replace(target, pattern, value, -1)
        }
    }
    url := fmt.Sprintf("%s://%s%s", config.Protocol, config.Host, target)
    //fmt.Println(url)
    resp, _ := http.Get(url)
    if resp == nil {
	fmt.Println("%s get nil response", url)
        return nil
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

func Post(config Config, forward Forwarding, context map[string]string, data_str string) []byte {
    target := forward.Target
    if data_str == "" {
        data_str = forward.Data
    }
    if context != nil {
        for param, value := range context {
            pattern := "$" + param
	    target = strings.Replace(target, pattern, value, -1)
	    data_str = strings.Replace(data_str, pattern, value, -1)
	}
    }
    data := []byte(data_str)
    url := fmt.Sprintf("%s://%s%s", config.Protocol, config.Host, target)
    resp, _ := http.Post(url, forward.ContentType, bytes.NewBuffer(data))
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

func Rpc(config Config, forward Forwarding, context map[string]string, data []byte) []byte {
    target := forward.Target
    host := forward.Host
    if context != nil {
        for param, value := range context {
            pattern := "$" + param
	    target = strings.Replace(target, pattern, value, -1)
	    host = strings.Replace(host, pattern, value, -1)
	}
    }
    url := fmt.Sprintf("%s://%s%s", config.Protocol, host, target)
    fmt.Println(url)
    fmt.Println(data)
    resp, _ := http.Post(url, forward.ContentType, bytes.NewBuffer(data))
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

func loadConfig(key string, fileName string) {
    var config Config

    jsonFile, err := os.Open(fileName)
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal([]byte(byteValue), &config)

    config.Indices = make(map[string] *Forwarding)
    for i, forwarding := range config.Forwardings {
        config.Indices[forwarding.Source] = &config.Forwardings[i]
    }

    Configurations[key] = config
}

func Initialize(configDir string) {
    var keys = []string{
        "tezos",
	"eos",
	"gxchain",
        "../action"}

    Configurations = make(map[string]Config)
    for _, key := range keys {
        loadConfig(key, configDir + key + ".json")
    }
}
