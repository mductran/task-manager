package taskmanager

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"sort"
)

const (
	CPU     = 1
	RAM     = 2
	DISK    = 3
	NETWORK = 4
	GPU     = 5
)

type Process struct {
	name          string
	pid           uint32
	parent_pid    uint32
	children      []Process
	cpu_usage     float32
	ram_usage     uint32
	disk_usage    float32
	network_usage float32
	// gpu_usage     float32
}

type SystemInfo struct {
	cpu_name       string
	cpu_freq       uint16
	cpu_temp       int8
	cpu_core_count uint8
	process_count  uint32
	uptime         float64
	mem_total      uint32
	processes      []Process
}

func List() ([]Process, error) {
	return list()
}

func Search(processes *[]Process, target interface{}) ([]Process, error) {
	switch v := target.(type) {
	case string:
		return searchByName(processes, v)
	case uint32:
		return searchByPid(processes, v)
	default:
		// something's going wrong here
		return []Process{}, errors.New("invalid type, need to be string or uint32")
	}
}

func Sort(processes []Process, target interface{}) error {
	switch target {
	case "CPU":
		sort.SliceStable(processes, func(i, j int) bool {
			return processes[i].cpu_usage < processes[j].cpu_usage
		})
		return nil
	case "RAM":
		sort.SliceStable(processes, func(i, j int) bool {
			return float32(processes[i].ram_usage) < float32(processes[j].ram_usage)
		})
		return nil
	case "DISK":
		sort.SliceStable(processes, func(i, j int) bool {
			return processes[i].disk_usage < processes[j].disk_usage
		})
		return nil
	case "NETWORK":
		sort.SliceStable(processes, func(i, j int) bool {
			return processes[i].network_usage < processes[j].network_usage
		})
		return nil
	default:
		return errors.New("wrong sort flag")
	}
}

func Start(executablePath string) (uint32, error) {
	executablePath = filepath.Clean(executablePath)
	currentDir, err := os.Getwd()
	if err != nil {
		return 0, err
	}
	return start(path.Join(currentDir, executablePath))
}

func Stop(targetPid uint32) (bool, error) {
	return stop(targetPid)
}

func Suspend(pid uint32) (bool, error) {
	return suspend(pid)
}

func Resume(pid uint32) (bool, error) {
	return resume(pid)
}
