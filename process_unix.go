//go:build unix

package taskmanager

import (
	"io"
	"os"
	"strconv"
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
