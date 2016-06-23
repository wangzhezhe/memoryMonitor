package alert

import (
	"testing"
	"time"
)

func TestTrigger(t *testing.T) {
	err := Trigger("nodeID", "containerID", 10)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 5)
	err = Resolve("nodeID", "containerID", 10)
	if err != nil {
		t.Error(err)
	}
}
