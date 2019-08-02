package consul

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

type ConsulBackup struct {
	Type      string
	Format    string
	Exclude   []string
	Host      string
	StartDate time.Time
}

func NewConsulBackup(m map[string]string) *ConsulBackup {

	return &ConsulBackup{
		StartDate: time.Now(),
	}
}

// Backup 备份
func (obj *ConsulBackup) Backup() {
	consulSDK := NewConsulAPI(obj.Host)
	cofigArray := consulSDK.Keys()
	// for _, v := range cofigArray {
	// 	fmt.Println(v)
	// }
	var fillerConfigArray []string
	if len(obj.Exclude) > 0 {
		for _, originConfig := range cofigArray {
			exclude := false
			for _, excludeRule := range obj.Exclude {
				re := regexp.MustCompile(excludeRule)
				exclude = re.Match([]byte(originConfig))
			}
			if !exclude && originConfig[len(originConfig)-1] != '/' {
				// fmt.Println(originConfig)
				fillerConfigArray = append(fillerConfigArray, originConfig)
			}
		}
	} else {
		for _, originConfig := range cofigArray {
			if originConfig[len(originConfig)-1] != '/' {
				fillerConfigArray = append(fillerConfigArray, originConfig)
			}
		}
	}
	for _, v := range fillerConfigArray {
		obj.SaveToLocal(v)
		// fmt.Println(v)
	}

}

// SaveToLocal 保存到本地,先这样实现了.
// fix me:最好定义一个底层可扩展的接口
func (obj *ConsulBackup) SaveToLocal(consulKey string) {
	today := obj.StartDate.Format("2006-01-02")
	now := obj.StartDate.Format("150405")
	consulSDK := NewConsulAPI(obj.Host)
	rawValue := consulSDK.Key(consulKey)
	consulURL, err := url.Parse(obj.Host)
	if err != nil {
		log.Fatal(err)
		return
	}
	filePath := path.Join("file", "consul", consulURL.Host, today, now, consulKey)
	parentDir := filepath.Dir(filePath)
	// fmt.Printf("filePath: %s \n", filePath)
	// fmt.Printf("parentDir: %s \n", parentDir)
	_, err = os.Stat(parentDir)
	if err != nil {
		os.MkdirAll(parentDir, os.ModePerm)
	}
	//全部只读
	err = ioutil.WriteFile(filePath, []byte(rawValue), 0444)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// CleanOld 危险接口,慎重使用
func (obj *ConsulBackup) CleanOld(before time.Time) {
	consulURL, err := url.Parse(obj.Host)
	if err != nil {
		log.Fatal(err)
		return
	}
	backupDir := path.Join("file", "consul", consulURL.Host)
	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, f := range files {
		// fmt.Println(f.Name())
		os.RemoveAll(path.Join(backupDir, f.Name()))

	}
}
