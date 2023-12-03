//go:build windows

package taskmanager

import (
	"fmt"
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

func splitAtCommas(s string) []string {
	res := []string{}
	var beg int
	var inString bool

	for i := 0; i < len(s); i++ {
		if s[i] == ',' && !inString {
			res = append(res, s[beg:i])
			beg = i + 1
		} else if s[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}

func getProcess(line string) (Process, error) {
	elements := splitAtCommas(line)

	pidString := clearString(elements[1])
	ramString := clearString(elements[4])
	ramString = ramString[:len(ramString)-2]

	fmt.Println(ramString)

	pid, _ := strconv.ParseUint(pidString, 10, 16)
	ram, _ := strconv.ParseUint(ramString, 10, 16)

	pw := Process{
		name:      elements[0],
		pid:       uint32(pid),
		ram_usage: uint32(ram),
	}
	return pw, nil
}

func parse(processes string) ([]Process, error) {
	pSlice := strings.Split(processes, "\n")

	var out []Process
	for _, p := range pSlice {
		if len(p) > 0 {
			process, err := getProcess(p)
			if err == nil {
				out = append(out, process)
			}
		}
	}

	return out, nil
}

func list() ([]Process, error) {
	cmd := exec.Command("tasklist", "/fo", "csv", "/nh")

	out, err := cmd.Output()
	if err != nil {
		return []Process{}, nil
	}

	processes, err := parse(string(out))
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
	// cmd := exec.Command("powershell", "-Command", "Start-Process", path)
	cmd := exec.Command(path)
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
		fmt.Println("handle: ", handle)
		return false, err2
	}

	return true, nil
}
