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
	// Type      string
	// Format    string
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

// SaveToLocal 保存到本地,先这样实现了.
// fix me:最好定义一个底层可扩展的接口
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
	filePath := path.Join(obj.PrefixPath, consulURL.Host, today, now, consulKey)
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

// CleanOld 危险接口,慎重使用
func (obj *ConsulBackup) CleanOld(cron string) {
	now := time.Now()
	crons := strings.Split(cron, " ")
	interval, err := strconv.ParseInt(crons[0], 10, 64)
	if err != nil {
		fmt.Errorf("cron format error.")
		os.Exit(-1)
		return
	}
	deadline := now
	switch crons[1] {
	case "d":
		deadline = deadline.AddDate(0, 0, -int(interval))
		break
	case "h":
		deadline = deadline.Add(time.Hour * -time.Duration(interval))
		break
	case "m":
		deadline = deadline.Add(time.Minute * -time.Duration(interval))
		break
	case "s":
		deadline = deadline.Add(time.Second * -time.Duration(interval))
		break
	}
	consulURL, err := url.Parse(obj.Host)
	if err != nil {
		// log.Fatal(err)
		return
	}
	// 遍历清除文件
	backupDir := path.Join(obj.PrefixPath, consulURL.Host)
	location, _ := time.LoadLocation("Local")
	deadlineDate := time.Date(deadline.Year(), deadline.Month(), deadline.Day(), 0, 0, 0, 0, location)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	// fmt.Printf("today: %v \n", today)
	deadlineAfterToday := deadline.After(today)
	// fmt.Printf("deadlineAfterToday: %v \n", deadlineAfterToday)
	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		// log.Fatal(err)
		return
	}
	for _, f := range files {
		//先删除前几天的文件,再删除今天之内的文件
		dayDir := f.Name()
		date, err := time.ParseInLocation(DAY_DIR_FORMAT, dayDir, location)
		if err != nil {
			// 存在其他类型的隐藏文件,所以这里 continue
			// fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}
		// fmt.Printf("date: %v \n", date)
		// fmt.Printf("today.Equal(date): %v \n", today.Equal(date))
		if today.Equal(date) && deadlineAfterToday {
			//清理当天文件夹
			todayFiles, _ := ioutil.ReadDir(path.Join(backupDir, f.Name()))
			for _, todayFile := range todayFiles {
				date, _ := time.ParseInLocation(DATE_DIR_FORMAT, todayFile.Name(), location)
				thatDate := time.Date(now.Year(), now.Month(), now.Day(), date.Hour(), date.Minute(), date.Second(), 0, location)
				// fmt.Printf("thatDate: %v \n", thatDate)
				if thatDate.Before(deadline) {
					deleteDir := path.Join(backupDir, f.Name(), todayFile.Name())
					fmt.Printf("delete Today dir: %v \n", deleteDir)
					os.RemoveAll(deleteDir)
				}
			}
		} else if date.Before(deadlineDate) {
			deleteDir := path.Join(backupDir, f.Name())
			fmt.Printf("delete old dir: %v \n", deleteDir)
			os.RemoveAll(deleteDir)
		}
	}
}
