package taskmanager

import (
	"bytes"
	"errors"
	"os/exec"
)

type ProcessWindows struct {
	Process
}

func parseProcessesWindows(processes string) ([]Process, error) {
	return nil, nil
}

func (p *ProcessWindows) List() ([]Process, error) {
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

	processes, err := parseProcessesWindows(out)
	if err != nil {
		return []Process{}, err
	}

	return processes, nil
}
