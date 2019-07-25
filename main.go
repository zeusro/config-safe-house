package main

import (
	"fmt"

	"sync"

	"github.com/jasonlvhit/gocron"
)

func main() {
	var waitGroup sync.WaitGroup
	fmt.Println(" config-safe-house is running.")
	waitGroup.Add(1)
	go func() {
		s1 := gocron.NewScheduler()
		s1.Every(1).Hours().Do(task)
		<-s1.Start()
		waitGroup.Done()
	}()
	waitGroup.Add(1)
	go func() {

		s2 := gocron.NewScheduler()
		s2.Every(1).Hours().Do(task2)
		<-s2.Start()
		waitGroup.Done()
	}()
	waitGroup.Wait()
	fmt.Println(" Exit. ")
}

func task() {
	fmt.Println("I am runnning task1.")
}

func task2() {
	fmt.Println("I am runnning task2.")
}
