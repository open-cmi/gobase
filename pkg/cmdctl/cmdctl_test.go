package cmdctl

import (
	"testing"
)

func TestNew(t *testing.T) {
	manager := NewManager()
	var conf Config = Config{
		Name:       "test",
		ExecStart:  "echo 1",
		RestartSec: 1,
		StopSignal: 2,
	}

	err := manager.AddProcess(&conf)
	if err != nil {
		t.Errorf("test: add process 1 failed")
	}

	err = manager.StartProcess("test")
	if err != nil {
		t.Errorf("test: start process failed")
	}

	err = manager.StopProcess("test")
	if err != nil {
		t.Errorf("test: stop process failed")
	}
	err = manager.DelProcess("test")
	if err != nil {
		t.Errorf("test: delete process failed")
	}
}
