package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/memoryMonitor/memory"
)

var (
	DefaultInterval = 60
	DefaultAlert    = 80
	NodeIP          string
)

func main() {

	flag.IntVar(&DefaultInterval, "interval", 60, "the interval to check the dir (s)")
	flag.IntVar(&DefaultAlert, "percentage", 80, "the percentage threshold to alert (rss+cache/total%)")
	flag.StringVar(&NodeIP, "nodeip", "", "the nodeip of local host")
	flag.Parse()
	if NodeIP == "" {
		fmt.Println("please input nodeip")
		return
	}
	//NemMemoryManager(interval int, alertThreshold int, nodeIP string)
	MemoManager, err := memory.NemMemoryManager(DefaultInterval, DefaultAlert, NodeIP)
	if err != nil {
		os.Exit(1)
	}
	for {
		MemoManager.Run()
	}

	fmt.Println("main exit for unknown reason")

}
