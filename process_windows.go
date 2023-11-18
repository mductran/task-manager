//go:build windows

package taskmanager

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

func clearString(str string) string {
	return regexp.MustCompile(`[^\p{L}\p{N} ]+`).ReplaceAllString(str, "")
}

func getProcess(line string) (Process, error) {
	line = clearString(line)
	elements := strings.SplitAfter(line, ",")

	pid, _ := strconv.ParseUint(elements[1], 10, 16)
	ram, _ := strconv.ParseUint(elements[4], 10, 16)
	// sess, _ := strconv.ParseUint(elements[3], 10, 16)
	pw := Process{
		name:      elements[0],
		pid:       uint32(pid),
		ram_usage: uint32(ram),
	}
	return pw, nil
}

func parse(processes string) ([]Process, error) {
	scanner := bufio.NewScanner(strings.NewReader(processes))
	i := 1
	var out []Process
	for scanner.Scan() {
		newProcess, err := getProcess(scanner.Text())
		if err == nil {
			out = append(out, newProcess)
		}
		i += 1
	}
	return out, nil
}

func list() ([]Process, error) {
	cmd := exec.Command("tasklist")

	var outBuffer, errBuffer bytes.Buffer
	if err := cmd.Run(); err != nil {
		return []Process{}, err
	}
	out := string(outBuffer.String())

	if errBuffer.Len() != 0 {
		err := errBuffer.String()
		return []Process{}, errors.New(err)
	}

	processes, err := parse(out)
	if err != nil {
		return []Process{}, err
	}

	return processes, nil
}

func searchByName(processes *[]Process, target string) ([]Process, error) {
	var out []Process
	for _, process := range *processes {
		if process.name == target {
			out = append(out, process)
		}
	}

	return out, nil
}

func searchByPid(processes *[]Process, target uint32) ([]Process, error) {
	var out []Process
	for _, process := range *processes {
		if process.pid == target {
			out = append(out, process)
		}
	}

	return out, nil
}

func start(path string) (uint32, error) {
	cmd := exec.Command("powershell", "-Command", "Start-Process", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return uint32(cmd.Process.Pid), nil
}

func stop(pid uint32) (bool, error) {
	cmd := exec.Command("taskkill", "/F", "/PID", strconv.FormatUint(uint64(pid), 10))
	if err := cmd.Run(); err != nil {
		return false, err
	}

	return true, nil
}

// Use win32 API to suspend processes
func suspend(pid uint32) (bool, error) {
	processHandle, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, pid)
	if err != nil {
		return false, err
	}
	defer windows.CloseHandle(processHandle)

	// https://www.vbforums.com/showthread.php?810801-So-here-s-how-to-suspend-a-process-via-Windows-API
	ntSuspendProcess := windows.NewLazySystemDLL("ntdll.dll").NewProc("NtSuspendProcess")
	r1, _, err := ntSuspendProcess.Call(uintptr(processHandle))
	if r1 != 0 {
		return false, err
	}

	return true, nil
}

func resume(pid uint32) (bool, error) {
	// PROCESS ACCESS RIGHT: allow all access rights to a process: CREATE, DUP, QUERY, SET, SUSPEND, TERMINATE, READ, WRITE, SYNC
	// https://learn.microsoft.com/en-us/windows/win32/procthread/process-security-and-access-rights
	var PROCESS_ALL_ACCESS = 0x1F0FFF

	// process is considered suspended if all threads it possesses are suspended
	// to resume the process, must resume all the threads it possesses
	var kernel32 = syscall.MustLoadDLL("kernel32.dll")
	openProcess := kernel32.MustFindProc("OpenProcess")
	resumeThread := kernel32.MustFindProc("ResumeThread")

	r1, _, err1 := openProcess.Call(uintptr(PROCESS_ALL_ACCESS), uintptr(1), uintptr(pid))
	if r1 == 0 {
		return false, err1
	}
	handle := syscall.Handle(r1)

	// resume the process
	r2, _, err2 := resumeThread.Call(uintptr(handle))
	if r2 == 0xFFFFFFFF {
		return false, err2
	}

	return true, nil
}
