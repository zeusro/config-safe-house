package consul

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
)

const (
	DAY_DIR_FORMAT  = "2006-01-02"
	DATE_DIR_FORMAT = "150405"
)

type ConsulBackup struct {
	Name       string
	Exclude    []string
	Host       string
	StartDate  time.Time
	PrefixPath string
}

func NewConsulBackup() *ConsulBackup {
	return &ConsulBackup{
		StartDate: time.Now(),
	}
}

func (obj *ConsulBackup) BackupByInterval(cron string) {
	s1 := gocron.NewScheduler()
	crons := strings.Split(cron, " ")
	interval, err := strconv.ParseUint(crons[0], 10, 64)
	if err != nil {
		fmt.Errorf("cron format error.")
		os.Exit(-1)
		return
	}
	switch crons[1] {
	case "d":
		s1.Every(interval).Days().Do(func() {
			obj.Backup()
		})
		break
	case "h":
		s1.Every(interval).Hours().Do(func() {
			obj.Backup()
		})
		break
	case "m":
		s1.Every(interval).Minutes().Do(func() {
			obj.Backup()
		})
		break
	case "s":
		s1.Every(interval).Seconds().Do(func() {
			obj.Backup()
		})
		break
	default:
		fmt.Errorf("cron format error.")
		os.Exit(-1)
		return
	}
	<-s1.Start()
}

// Backup 备份
func (obj *ConsulBackup) Backup() {
	obj.StartDate = time.Now()
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

// SaveToLocal 保存到本地磁盘
func (obj *ConsulBackup) SaveToLocal(consulKey string) {
	today := obj.StartDate.Format(DAY_DIR_FORMAT)
	now := obj.StartDate.Format(DATE_DIR_FORMAT)
	consulSDK := NewConsulAPI(obj.Host)
	rawValue := consulSDK.Key(consulKey)
	consulURL, err := url.Parse(obj.Host)
	if err != nil {
		log.Fatal(err)
		return
	}
	filePath := path.Join(obj.PrefixPath, obj.Name+"_"+consulURL.Host, today, now, consulKey)
	parentDir := filepath.Dir(filePath)
	// fmt.Printf("filePath: %s \n", filePath)
	// fmt.Printf("parentDir: %s \n", parentDir)
	_, err = os.Stat(parentDir)
	if err != nil {
		os.MkdirAll(parentDir, os.ModePerm)
	}
	//全部只读
	err = ioutil.WriteFile(filePath, []byte(rawValue), 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// replaceText 配置替换字典
func replaceText(text string) string {
	m := make(map[string]string)
	// m[""] = ""
	// m[""] = ""
	for k, v := range m {
		if strings.Contains(text, k) {
			fmt.Println("包含 ", k)
			text = strings.ReplaceAll(text, k, v)
		}
	}
	return text
}
