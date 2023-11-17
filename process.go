package taskmanager

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
}

func List() ([]Process, error) {
	list()
	return nil, nil
}

func Search() ([]Process, error) {
	return nil, nil
}

func Sort() ([]Process, error) {
	return nil, nil
}

func Start() (bool, error) {
	return false, nil
}

func Stop() (bool, error) {
	return false, nil
}

func Pause() (bool, error) {
	return false, nil
}

func Resume() (bool, error) {
	return false, nil
}
