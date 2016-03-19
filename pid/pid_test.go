package pid

import (
	"os"
	"testing"
)

func TestPid(t *testing.T) {
	pidFile := "/tmp/test.pid"
	err := CreatePidFile(&pidFile)
	defer os.Remove(pidFile)
	assert(t, err == nil, "err create pid file")

	assert(t, Exists(pidFile), "can't create pid file")
	RemovePidFile(&pidFile)
	assert(t, !Exists(pidFile), "can't delete pid file")
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func assert(t *testing.T, ok bool, msg string) {
	if !ok {
		t.Error(msg)
	}
}
