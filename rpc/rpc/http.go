package rpc

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
    fmt.Println(url)
    resp, _ := http.Get(url)
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

func Post(config Config, forward Forwarding, context map[string]string) []byte {
    target := forward.Target
    data_str := forward.Data
    if context != nil {
        for param, value := range context {
            pattern := "$" + param
	    target = strings.Replace(target, pattern, value, -1)
	    data_str = strings.Replace(data_str, pattern, value, -1)
	}
    }
    url := fmt.Sprintf("%s://%s%s", config.Protocol, config.Host, target)
    data := []byte(data_str)
    resp, _ := http.Post(url, forward.ContentType, bytes.NewBuffer(data))
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

/*
func (config *Config) readConfigFile(filename string) {
    jsonFile, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal([]byte(byteValue), config)

    config.Indices = make(map[string] *URLType)
    for i, forwarding := range config.Forwardings {
        config.Indices[forwarding.source] = &config.Forwardings[i]
    }
}
*/

func loadConfig(key string, fileName string) {
    var config Config

    jsonFile, err := os.Open(fileName)
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal([]byte(byteValue), config)

    config.Indices = make(map[string] *Forwarding)
    for i, forwarding := range config.Forwardings {
        config.Indices[forwarding.Source] = &config.Forwardings[i]
    }

    Configurations[key] = config
}

func Initialize() {
    configDir := "../../config/query/"
    var keys = []string{
        "tezos",
	"eos",
	"gxchain"}

    Configurations = make(map[string]Config)
    for _, key := range keys {
        loadConfig(key, configDir + key + ".json")
    }
}
