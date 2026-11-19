//go:build unix

package executor

import (
	"os"
	"os/exec"
	"syscall"
)

// configureProcessGroup puts the command in its own process group so that
// Terraform's provider plugins are stopped together with the parent process.
func configureProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcessGroup terminates the command and every process it spawned.
func killProcessGroup(process *os.Process) {
	if process == nil {
		return
	}
	_ = syscall.Kill(-process.Pid, syscall.SIGKILL)
	_ = process.Kill()
}

func runningAsRoot() bool {
	return os.Geteuid() == 0
}
