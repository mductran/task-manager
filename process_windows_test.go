//go:build windows

package taskmanager

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
	"testing"
)

func TestStartProcess(t *testing.T) {
	pid, err := Start("./example.exe")
	if err != nil {
		t.Error(err)
	}
	if pid == 0 {
		t.Errorf("could not start example process")
	}

	// time.Sleep(2 * time.Second)

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
}
