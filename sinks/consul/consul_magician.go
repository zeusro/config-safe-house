package consul

import (
	"fmt"
	"strings"
	"sync"
)

type ConsulMagician struct {
	Host string
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
	newValue := replaceText(value)
	if strings.EqualFold(value, newValue) {
		return
	}
	fmt.Printf("deal wit key: %v \n", key)
	consulSDK.UpdateKey(key, newValue)
}
