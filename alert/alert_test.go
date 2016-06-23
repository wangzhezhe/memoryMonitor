package alert

import (
	"testing"
)

func TestMemoryAlert(t *testing.T) {
	nodeIP := "test"
	containerID := "test"
	memoryPercentage := 60
	err := MemoryAlert(nodeIP, containerID, memoryPercentage)
	if err != nil {
		t.Error(err)
	}

}
