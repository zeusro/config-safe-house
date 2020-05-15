package main

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/jasonlvhit/gocron"
	"github.com/zeusro/config-safe-house/model"
	"github.com/zeusro/config-safe-house/sinks/consul"
)

func main() {
	var waitGroup sync.WaitGroup
	fmt.Println("config-safe-house is running.")
	configPath := path.Join("config.yaml")
	config, _ := model.ParseFromPath(configPath)
	if config == nil {
		configPath = path.Join("config-default.yaml")
		config, _ = model.ParseFromPath(configPath)
	}
	if config == nil {
		fmt.Printf("Load config fail, please check if config.yaml/config-default.yaml exist. \n path: %s \n", configPath)
		os.Exit(-1)
		return
	}
	fmt.Printf("config: %#v \n", config)
	for _, consulConfig := range config.Consul {
		if err := consulConfig.CheckSelf(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(-1)
			return
		}
		waitGroup.Add(1)
		//对象的内容(consulConfig)有变化,但是指针是不变的,所以这里要传对象的复制
		go func(consulInfo model.ConsulConfig) {
			fmt.Printf("backup job : %s ;\n", consulInfo.InstanceURL)
			backupConfigs := consulInfo.Backup
			for _, backupConfig := range backupConfigs {
				// waitGroup.Add(1)
				// 同理,对象的内容有变化,但是指针是不变的,所以这里要传对象的复制
				go func(consulURL string, backupStrategy model.BackupStrategy) {
					backupJob := consul.NewConsulBackup()
					backupJob.Name = consulInfo.InstanceName
					backupJob.Exclude = backupStrategy.File.Exclude
					backupJob.Host = consulURL
					backupJob.PrefixPath = backupStrategy.File.Path
					cleanPolicy := backupStrategy.File.CleanPolicy
					if len(cleanPolicy) > 0 {
						//清除旧的备份文件
						// waitGroup.Add(1)
						go func() {
							s1 := gocron.NewScheduler()
							s1.Every(1).Minutes().Do(func() {
								killer := consul.ConsulKiller{
									Name:       consulInfo.InstanceName,
									Host:       consulURL,
									PrefixPath: backupStrategy.File.Path,
								}
								killer.CleanOld(cleanPolicy)
							})
							<-s1.Start()
							// waitGroup.Done()
						}()
					}
					//解析表达式
					cron := backupStrategy.File.Cron
					backupJob.BackupByInterval(cron)
					// waitGroup.Done()
				}(consulInfo.InstanceURL, backupConfig)
			}
			// waitGroup.Done()
		}(consulConfig)
	}
	// wait till world end
	waitGroup.Wait()
	fmt.Println("Exit. ")
}
