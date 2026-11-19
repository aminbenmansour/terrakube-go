//go:build unix

package executor

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestRunCommandKillsProcessGroup proves that cancelling a job stops the
// Terraform process and the provider children it spawned, not just the direct
// child.
func TestRunCommandKillsProcessGroup(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "spawn.sh")
	body := "#!/bin/sh\n" +
		"echo $$ > parent.pid\n" +
		"sleep 30 &\n" +
		"echo $! > child.pid\n" +
		"wait\n"
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	e := New(Deps{})
	sink := newLogSink(context.Background(), nil, nil, 1, "step-1", nil)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- e.runCommand(ctx, dir, nil, sink, script) }()

	parentPID := waitForPID(t, filepath.Join(dir, "parent.pid"))
	childPID := waitForPID(t, filepath.Join(dir, "child.pid"))
	cancel()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("runCommand did not return after cancellation")
	}
	for _, pid := range []int{parentPID, childPID} {
		if !waitForProcessExit(pid, 3*time.Second) {
			t.Fatalf("process %d survived cancellation", pid)
		}
	}
}

func waitForPID(t *testing.T, path string) int {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		data, err := os.ReadFile(path)
		if err == nil {
			pid, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
			if parseErr == nil && pid > 0 {
				return pid
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("pid file %s was not written: %v", path, err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func waitForProcessExit(pid int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for {
		if err := syscall.Kill(pid, 0); errors.Is(err, syscall.ESRCH) {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(10 * time.Millisecond)
	}
}
