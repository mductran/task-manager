//go:build windows

package taskmanager

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
	"testing"
)

// TODO: remove this comment
// start process -> check started pid match -> stop process -> check if pid still persist
func TestStartStopProcess(t *testing.T) {

	// Start
	pid, err := Start("./example.exe")
	if err != nil {
		t.Error(err)
	}
	if pid == 0 {
		t.Errorf("could not start example process")
	}

	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Get-Process -Id %d", pid))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	if err = cmd.Run(); err != nil {
		t.Error("Fail to run Get-Process: ", err)
	}
	if len(errb.String()) > 0 {
		t.Errorf("Fail to run Get-Grocess: %v", errb.String())
	}
	if len(outb.String()) < 1 {
		t.Errorf("Cannot find process with PID: %v", pid)
	}
	fmt.Println(outb.String())

	// Suspend

	// Resume

	// Stop
	killed, err := Stop(pid)
	if err != nil {
		t.Errorf("While terminating process with pid %v: ", err)
	}
	if err = cmd.Run(); err == nil {
		t.Error("Process was not killed: ", err)
	}

	fmt.Println("Process is stopped: ", killed)
}
