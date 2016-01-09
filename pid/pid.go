package pid

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func CreatePidFile(pidFile string) error {
	if pidString, err := ioutil.ReadFile(pidFile); err == nil {
		pid, err := strconv.Atoi(string(pidString))
		if err == nil {
			if _, err := os.Stat(fmt.Sprintf("/proc/%d/", pid)); err == nil {
				return fmt.Errorf("pid file found, ensure stns  is not running or delete %s", pidFile)
			}
		}
	}

	file, err := os.Create(pidFile)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = fmt.Fprintf(file, "%d", os.Getpid())
	return err
}

func RemovePidFile(pidFile string) {
	if err := os.Remove(pidFile); err != nil {
		log.Fatal("Error removing %s: %s", pidFile, err)
	}
}
