package consul

import (
	"fmt"
	"strings"
	"sync"
)

type ConsulMagician struct {
	Host   string
	DryRun bool
}

// ReplaceAllKeys 替换所有 key
func (obj *ConsulMagician) ReplaceAllKeys() {
	consulSDK := NewConsulAPI(obj.Host)
	cofigArray := consulSDK.Keys()
	var wg sync.WaitGroup
	for _, key := range cofigArray {
		wg.Add(1)
		go func(key string) {
			obj.ReplacelKey(key)
			wg.Done()
		}(key)
	}
	wg.Wait()
}

func (obj *ConsulMagician) ReplacelKey(key string) {
	consulSDK := NewConsulAPI(obj.Host)
	value := consulSDK.Key(key)
	newValue := obj.replaceText(value)
	if strings.EqualFold(value, newValue) {
		return
	}
	fmt.Printf(" deal wit key: %v \n", key)
	if !obj.DryRun {
		consulSDK.UpdateKey(key, newValue)
	}
}

// replaceText 配置替换字典
func (obj *ConsulMagician) replaceText(text string) string {
	m := make(map[string]string)
	// m["old-value"] = "new-value"
	// m[""] = ""
	// m[""] = ""
	for k, v := range m {
		if strings.Contains(text, k) {
			fmt.Print("contain " + k)
			text = strings.ReplaceAll(text, k, v)
		}
	}
	return text
}
