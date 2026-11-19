//go:build !unix

package executor

import (
	"os"
	"os/exec"
)

func configureProcessGroup(*exec.Cmd) {}

func killProcessGroup(process *os.Process) {
	if process != nil {
		_ = process.Kill()
	}
}

func runningAsRoot() bool {
	return false
}
