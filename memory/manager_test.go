package memory

import (
	"testing"
)

func TestGetContainerMemCapacity(t *testing.T) {
	path := "/sys/fs/cgroup/memory/docker/0f3b64bc485bf370e68d72851d5b7d3eefacb93784644575bdeccce6edc43325"
	value, err := getContainerMemCapacity(path)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("value: ", value)
}

func TestRun(t *testing.T) {
	//t.SkipNow()
	var interval = 1
	manager, err := NemMemoryManarer(interval)
	if err != nil {
		t.Error(err)
	}

	t.Log(manager)

	manager.Run()
}
