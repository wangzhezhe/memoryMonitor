package memory

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/memoryMonitor/alert"
)

var (
	alertRecord = make(map[string]int)
)

var (
	DefaultCgroupDir        = "/sys/fs/cgroup"
	DefaultMemoInfoFile     = "/proc/meminfo"
	DefaultDockerReg        = regexp.MustCompile(`[a-zA-Z0-9-_+.]+:[a-fA-F0-9]+`)
	memoryCapacityRegexp    = regexp.MustCompile(`MemTotal:\s*([0-9]+) kB`)
	swapCapacityRegexp      = regexp.MustCompile(`SwapTotal:\s*([0-9]+) kB`)
	cgroupMemCapacityRegexp = regexp.MustCompile(`([0-9]+)`)
)

type MemoryManager struct {
	Interval       int //wait Interval second to get the memo info
	MemoTotal      int //Bytes
	AlertThreshold int //cache + rss %
	NodeIP         string
}

// copy from cadvisor
// parseCapacity matches a Regexp in a []byte, returning the resulting value in bytes.
// Assumes that the value matched by the Regexp is in KB.
func parseCapacity(b []byte, r *regexp.Regexp) (int, error) {
	matches := r.FindSubmatch(b)
	if len(matches) != 2 {
		return -1, fmt.Errorf("failed to match regexp in output: %q", string(b))
	}
	m, err := strconv.ParseInt(string(matches[1]), 10, 64)
	if err != nil {
		return -1, err
	}

	// Convert to bytes.
	return int(m * 1024), err
}

func getContainerMemCapacity(path string) (int, error) {
	fileName := path + "/" + "memory.usage_in_bytes"
	out, err := ioutil.ReadFile(fileName)
	if err != nil {
		return -1, err
	}
	matches := cgroupMemCapacityRegexp.FindSubmatch(out)
	value, err := strconv.ParseInt(string(matches[0]), 10, 64)
	if err != nil {
		return -1, err
	}
	//log.Printf("the mem capacity for %s is %d", path, value)
	return int(value), nil
}

func getTotalMem() (int, error) {

	out, err := ioutil.ReadFile(DefaultMemoInfoFile)
	if err != nil {
		return -1, err
	}

	memoryCapacity, err := parseCapacity(out, memoryCapacityRegexp)
	if err != nil {
		return -1, err
	}

	return memoryCapacity, nil

}

func NemMemoryManager(interval int, alertThreshold int, nodeIP string) (*MemoryManager, error) {

	memTotal, err := getTotalMem()
	if err != nil {
		return nil, err
	}
	manager := &MemoryManager{
		Interval:       interval,
		MemoTotal:      memTotal,
		AlertThreshold: alertThreshold,
		NodeIP:         nodeIP,
	}
	return manager, nil
}

func (manager *MemoryManager) CheckMemCapacity() error {
	//traverse and check
	CgroupMemDir := DefaultCgroupDir + "/memory" + "/docker"
	dir, err := ioutil.ReadDir(CgroupMemDir)
	if err != nil {
		return nil
	}
	for _, fi := range dir {
		if fi.IsDir() {
			//TODO use regex to check the file name here
			containerID := fi.Name()
			path := CgroupMemDir + "/" + containerID
			containerMemValue, err := getContainerMemCapacity(path)
			if err != nil {
				log.Printf("failed to read the mem info from %s , the reason: %s", fi.Name(), err.Error())
				continue
			}
			if containerMemValue == -1 {
				return fmt.Errorf("the mem value for container %s is -1", fi.Name())
			}
			currentPercentage := float32(float32(containerMemValue)/float32(manager.MemoTotal)) * 100
			log.Printf("the percentage for container %s is %f", fi.Name(), currentPercentage)

			if int(currentPercentage) > manager.AlertThreshold {
				//TODO put the times into a variable
				if alertRecord[containerID] < 1 {
					alert.Trigger(manager.NodeIP, containerID, int(currentPercentage))
					alert.MemoryAlert(manager.NodeIP, containerID, int(currentPercentage))
				}

				alertRecord[containerID]++

			}
			if int(currentPercentage) < manager.AlertThreshold {
				alert.Resolve(manager.NodeIP, containerID, int(currentPercentage))

			}
		}
	}

	return nil

}

func (manager *MemoryManager) Run() {
	//modify the value in map to 0
	//TODO put the value into a variable
	recordTicker := time.NewTicker(time.Second * 600)
	go func() {
		for _ = range recordTicker.C {
			//create a new struct
			alertRecord = make(map[string]int)
			log.Println("-----------------------")
			log.Println("create a new record map")
		}
	}()

	//start to check the usage by internal
	ticker := time.NewTicker(time.Second * time.Duration(manager.Interval))

	for t := range ticker.C {
		fmt.Println("Tick at", t)
		manager.CheckMemCapacity()
	}
	//TODO add while to persist the exit accidentally
	fmt.Println("range exit for unknown reason")

}
