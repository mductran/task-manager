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
	ram_usage     float32
	disk_usage    float32
	network_usage float32
	// gpu_usage     float32
}

type ProcessInterface interface {
	Parse() Process

	List() ([]Process, error)

	Update()
}
