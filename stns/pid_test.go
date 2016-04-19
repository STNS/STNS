package stns

import (
	"os"
	"testing"

	"github.com/STNS/STNS/test"
)

func TestPid(t *testing.T) {
	pidFile := "/tmp/test.pid"
	err := createPidFile(pidFile)
	defer os.Remove(pidFile)
	test.Assert(t, err == nil, "err create pid file")

	test.Assert(t, Exists(pidFile), "can't create pid file")
	removePidFile(pidFile)
	test.Assert(t, !Exists(pidFile), "can't delete pid file")
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
