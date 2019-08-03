package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/jasonlvhit/gocron"
	"github.com/zeusro/config-safe-house/model"
	"github.com/zeusro/config-safe-house/sinks/consul"
)

func main() {
	var waitGroup sync.WaitGroup
	fmt.Println(" config-safe-house is running.")
	configPath := path.Join("config.yaml")
	config, _ := model.ParseFromPath(configPath)
	if config == nil {
		configPath = path.Join("config-default.yaml")
		config, _ = model.ParseFromPath(configPath)
	}
	if config == nil {
		fmt.Printf("Load config fail, please check if config.yaml/config-default.yaml exist")
		os.Exit(-1)
		return
	}
	fmt.Printf("%#v \n", config)
	for _, consulConfig := range config.Consul {
		if err := consulConfig.CheckSelf(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(-1)
			return
		}
		waitGroup.Add(1)
		go func() {
			fmt.Printf("backup job : %s ; \n ", consulConfig.InstanceURL)
			backupConfigs := consulConfig.Backup
			for _, backupConfig := range backupConfigs {
				waitGroup.Add(1)
				go func() {
					s1 := gocron.NewScheduler()
					//解析表达式
					cron := backupConfig.File.Cron
					crons := strings.Split(cron, " ")
					interval, err := strconv.ParseUint(crons[0], 10, 64)
					if err != nil {
						fmt.Errorf("cron format error.")
						os.Exit(-1)
						return
					}
					cleanPolicy := backupConfig.File.CleanPolicy
					if len(cleanPolicy) > 0 {
						//清除旧的备份文件
					}
					job := consul.NewConsulBackup()
					job.Exclude = backupConfig.File.Exclude
					job.Host = consulConfig.InstanceURL
					job.PrefixPath = backupConfig.File.Path
					switch crons[1] {
					case "d":
						s1.Every(interval).Days().Do(func() {
							job.Backup()
						})
						break
					case "h":
						s1.Every(interval).Hours().Do(func() {
							job.Backup()
						})
						break
					case "m":
						s1.Every(interval).Minutes().Do(func() {
							job.Backup()
						})
						break
					case "s":
						s1.Every(interval).Seconds().Do(func() {
							job.Backup()
						})
						break
					default:
						fmt.Errorf("cron format error.")
						os.Exit(-1)
						return
					}
					<-s1.Start()
					waitGroup.Done()
				}()
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
	fmt.Println(" Exit. ")
}

func cleanJob(date string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {

		s1 := gocron.NewScheduler()
		s1.Every(1).Minutes().Do(func() {
			// job.Backup()
		})
		wg.Done()
	}()

}
