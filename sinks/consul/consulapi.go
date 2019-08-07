package consul

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/zeusro/config-safe-house/util"
)

type ConsulAPI struct {
	Host string
}

const (
	// 遍历 key 接口
	URL_KEYS = "%s/v1/kv/?keys"
	// key 接口
	URL_KEY = "%s/v1/kv/%s?raw"
)

func NewConsulAPI(host string) *ConsulAPI {
	return &ConsulAPI{
		Host: host,
	}
}

func (obj *ConsulAPI) Keys() (cofigArray []string) {
	host := obj.Host
	configs, err := util.HttpGet(fmt.Sprintf(URL_KEYS, host))
	if err != nil {
		log.Fatal(err)
		return cofigArray
	}
	// var cofigArray []string
	err = json.Unmarshal([]byte(configs), &cofigArray)
	if err != nil {
		log.Fatal(err)
		return cofigArray
	}
	err = json.Unmarshal([]byte(configs), &cofigArray)
	if err != nil {
		log.Fatal(err)
		return cofigArray
	}
	return cofigArray
}

// Key https://www.consul.io/api/kv.html
func (obj *ConsulAPI) Key(key string) string {
	url := fmt.Sprintf(URL_KEY, obj.Host, key)
	// fmt.Println(url)
	response, err := util.HttpGet(url)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(response)
	return response
}

func (obj *ConsulAPI) UpdateKey(key, value string) {
	url := fmt.Sprintf(URL_KEY, obj.Host, key)
	_, err := util.HttpPut(url, value)
	if err != nil {
		fmt.Print(err)
		return
	}
	// fmt.Println(response)

}
