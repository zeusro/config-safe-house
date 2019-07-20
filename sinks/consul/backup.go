package consul

import (
	"encoding/json"
	"fmt"

	"github.com/zeusro/config-safe-house/util"
)

type ConsulBackup struct {
	Type    string
	Format  string
	exclude []string
	Host    string
}

func NewConsulBackup(m map[string]string) *ConsulBackup {

	return &ConsulBackup{}
}

func (obj *ConsulBackup) Backup() {
	host := obj.Host
	configs, err := util.HttpGet(fmt.Sprintf(URL_KEYS, host))
	if err != nil {
		fmt.Print(err)
		return
	}
	var cofigArray []string
	err = json.Unmarshal([]byte(configs), &cofigArray)
	if err != nil {
		fmt.Print(err)
		return
	}
 for _,v:=range cofigArray{
	 fmt.Println(v)
 }
}
