package consul

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type ConsulKiller struct {
	Name       string
	Host       string
	PrefixPath string
}

// CleanOld 危险接口,慎重使用
func (obj *ConsulKiller) CleanOld(cron string) {
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
	backupDir := path.Join(obj.PrefixPath, obj.Name+"_"+consulURL.Host)
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
