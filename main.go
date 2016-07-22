package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	installSignalHandler()

	var wg sync.WaitGroup
	wg.Add(1)
	go MemoManager.Run()
	wg.Wait()
	fmt.Println("main exit for unknown reason")

}

func installSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	// Block until a signal is received.
	go func() {
		sig := <-c
		log.Println("get the signal from os")
		log.Printf("Exiting given signal: %+v", sig)
		os.Exit(0)
	}()
}
