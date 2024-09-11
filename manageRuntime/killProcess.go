package manageRuntime

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func KillRunningInstances() {
	currentPID := os.Getpid()

	execPath, err := os.Executable()
	if err != nil {
		log.Println("Error getting executable path:", err)
		return
	}

	binaryName := filepath.Base(execPath)

	output, err := exec.Command("ps", "aux").Output()
	if err != nil {
		log.Println("Error getting process list:", err)
		return
	}

	processes := strings.Split(string(output), "\n")

	for _, process := range processes {
		if strings.Contains(process, binaryName) {
			fields := strings.Fields(process)
			if len(fields) < 2 {
				log.Println("Error malformed arg")
				continue
			}

			pid := fields[1]

			if pid == fmt.Sprint(currentPID) {
				continue
			}

			pidInt := stringToInt(pid)
			if err := syscall.Kill(pidInt, syscall.SIGKILL); err != nil {
				log.Println("Failed to kill process:", err)
			} else {
				log.Println("Killed process with PID:", pidInt)
			}
		}
	}
}

func stringToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Println("Failed to convert string to int:", err)
		return 0
	}
	return n
}
