//go:build unix

package taskmanager

import (
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func parse(pid int) (Process, error) {
	return Process{}, nil
}

func list() ([]Process, error) {

	dir, err := os.Open("/proc")
	if err != nil {
		return []Process{}, err
	}
	defer dir.Close()

	results := make([]int64, 0, 50)

	for {
		names, err := dir.Readdirnames(10)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		for _, name := range names {
			if name[0] < '0' || name[0] > '9' {
				continue
			}
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}
			results = append(results, pid)
		}
	}

	var processes []Process

	for result := range results {
		process, err := parseProcessesUnix(result)
		if err != nil {
			continue
		}
		processes = append(processes, process)
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
	cmd := exec.Command(path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return uint32(cmd.Process.Pid, nil)
}

func stop(pid uint32) (bool, error) {
	// TODO: sigkill v. sigterm
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return false, err
	}

	return true, nil
}

func suspend(pid uint32) (bool, error) {
	process, _ := os.FindProcess(pid)
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}
	err = process.Signal(syscall.SIGSTOP)
	if err != nil {
		return false, err
	}

	return true, nil
}

func resume(pid uint32) (bool, error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}
	err = process.Signal(syscall.SIGCONT)
	if err != nil {
		return false, err
	}

	return true, nil
}
