package taskmanager

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
)

type SystemInfo struct {
	cpu_hist        float32
	cpu_temp        float32
	gpu_hist        float32
	gpu_temp        float32
	ram_hist        float32
	uptime          float32
	processes_count uint32
}

type ResourceHistory struct {
	log []float32
}

func (r *ResourceHistory) LogData(val float32) {
	// push new data to front and remove last
	r.log = append([]float32{val}, r.log...)
}

// Sort Processes by fields
func SortProcesses(processes []Process, parents []Process, field int) {
	if field == CPU {
		sort.SliceStable(processes, func(i, j int) bool {
			return processes[i].cpu_usage < processes[j].cpu_usage
		})
	}
	if field == RAM {
		sort.SliceStable(processes, func(i, j int) bool {
			return processes[i].ram_usage < processes[j].ram_usage
		})
	}
	if field == DISK {
		sort.SliceStable(processes, func(i, j int) bool {
			return processes[i].disk_usage < processes[j].disk_usage
		})
	}
	if field == NETWORK {
		sort.SliceStable(processes, func(i, j int) bool {
			return processes[i].network_usage < processes[j].network_usage
		})
	}
	// if field == GPU {
	// 	sort.SliceStable(processes, func(i, j int) bool {
	// 		return processes[i].gpu_usage < processes[j].gpu_usage
	// 	})
	// }
}

func SearchProcessById() (Process, error) {
	return Process{}, nil
}

func SearchProcessByName() (Process, error) {
	return Process{}, nil
}

func getRuntime() string {
	return runtime.GOOS
}

func main() {
	fmt.Println(getRuntime())

	cmd := exec.Command("ls", "-lah")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdoutBuf.String()), string(stderrBuf.String())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
