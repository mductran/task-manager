//go:build windows

package taskmanager

import (
	"bytes"
	"fmt"
	"os/exec"
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

	getProcess := exec.Command("powershell", "-Command", fmt.Sprintf("Get-Process -Id %d", pid))

	var outb, errb bytes.Buffer
	getProcess.Stdout = &outb
	getProcess.Stderr = &errb
	if err = getProcess.Run(); err != nil {
		t.Error("Fail to run Get-Process: ", err)
	}
	if len(errb.String()) > 0 {
		t.Errorf("Fail to run Get-Grocess: %v", errb.String())
	}
	if len(outb.String()) < 1 {
		t.Errorf("Cannot find process with PID: %v", pid)
	}

	// Suspend
	_, err = Suspend(pid)
	if err != nil {
		t.Errorf("While suspending process with pid %v: ", err)
	}
	processStatus := exec.Command("powershell", "-Command", fmt.Sprintf("Get-Process -Id %d | select -expand Responding", pid))
	outb.Reset()
	errb.Reset()
	processStatus.Stdout = &outb
	processStatus.Stderr = &outb
	if err = processStatus.Run(); err != nil {
		t.Error("Fail to run Get-Process status: ", err)
	}
	if len(errb.String()) > 0 {
		t.Error("Fail to run Get-Process status: ", errb.String())
	}
	if outb.String() == "True" {
		t.Errorf("Process is not suspended")
	}

	// Resume
	_, err = Resume(pid)
	if err != nil {
		t.Errorf("While resuming process with pid: %v", err)
	}
	outb.Reset()
	errb.Reset()
	if err = processStatus.Run(); err != nil {
		t.Error("Fail to run Get-Process status: ", err)
	}
	if len(errb.String()) > 0 {
		t.Error("Fail to run Get-Process status: ", errb.String())
	}
	if outb.String() != "True" {
		t.Errorf("Process is not resumed")
	}

	// Stop
	_, err = Stop(pid)
	outb.Reset()
	errb.Reset()
	if err != nil {
		t.Errorf("While terminating process with pid %v: ", err)
	}
	if err = getProcess.Run(); err == nil {
		t.Error("Process was not killed: ", err)
	}
}
