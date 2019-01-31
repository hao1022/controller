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

type URLType struct {
    URL string `json:"url"`
    Source string `json:"source"`
    Method string `json:"method"`
    Data string `json:"data"`
    ContentType string `json:"content_type"`
    Params []string `json:"params"`
    Query []string `json:"query"`
}

type ConfigType struct {
    Root string `json:"root"`
    Endpoint string `json:"endpoint"`
    Protocol string `json:"protocol"`
    URLs []URLType `json:"URLs"`
    Indices map[string]*URLType
}

var Config map[string]ConfigType

func (config *ConfigType) ReadConfigFile(filename string) {
    jsonFile, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal([]byte(byteValue), config)

    config.Indices = make(map[string] *URLType)
    for i, url := range config.URLs {
        config.Indices[url.URL] = &config.URLs[i]
    }
}

func Get(config ConfigType, route URLType, context map[string]string) []byte {
    source := route.Source
    if context != nil {
        for param, value := range context {
            pattern := "$" + param
            source = strings.Replace(source, pattern, value, -1)
        }
    }
    url := fmt.Sprintf("%s://%s%s", config.Protocol, config.Endpoint, source)
    //fmt.Println(url)
    resp, _ := http.Get(url)
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

func Post(config ConfigType, route URLType, context map[string]string) []byte {
    source := route.Source
    data_str := route.Data
    if context != nil {
        for param, value := range context {
            pattern := "$" + param
	    source = strings.Replace(source, pattern, value, -1)
	    data_str = strings.Replace(data_str, pattern, value, -1)
	}
    }
    url := fmt.Sprintf("%s://%s%s", config.Protocol, config.Endpoint, source)
    data := []byte(data_str)
    resp, _ := http.Post(url, route.ContentType, bytes.NewBuffer(data))
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

func RegisterConfig(key string, config_file string) {
    var config ConfigType

    config.ReadConfigFile(config_file)
    Config[key] = config
}

func Initialize() {
    config_dir := "../../config/"
    var keys = []string{
        "rpc",   // RPC route
        "tezos",
	"eos",
	"gxchain"}

    Config = make(map[string]ConfigType)
    for _, key := range keys {
        RegisterConfig(key, config_dir + key + ".json")
    }
}
