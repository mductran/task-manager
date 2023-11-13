package taskmanager

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type ProcessWindows struct {
	Process
	session_name string
	session_num  uint32
}

func clearString(str string) string {
	return regexp.MustCompile(`[^\p{L}\p{N} ]+`).ReplaceAllString(str, "")
}

func process(line string) (ProcessWindows, error) {
	line = clearString(line)
	elements := strings.SplitAfter(line, ",")

	pid, _ := strconv.ParseUint(elements[1], 10, 16)
	ram, _ := strconv.ParseUint(elements[4], 10, 16)
	sess, _ := strconv.ParseUint(elements[3], 10, 16)
	pw := ProcessWindows{
		Process: Process{
			name:      elements[0],
			pid:       uint32(pid),
			ram_usage: uint32(ram),
		},
		session_name: elements[2],
		session_num:  uint32(sess),
	}
	return pw, nil
}

func parseProcessesWindows(processes string) ([]ProcessWindows, error) {
	scanner := bufio.NewScanner(strings.NewReader(processes))
	i := 1
	var out []ProcessWindows
	for scanner.Scan() {
		newProcess, err := process(scanner.Text())
		if err == nil {
			out = append(out, newProcess)
		}
		i += 1
	}
	return out, nil
}

func (p *ProcessWindows) List() ([]ProcessWindows, error) {
	cmd := exec.Command("tasklist")

	var outBuffer, errBuffer bytes.Buffer
	if err := cmd.Run(); err != nil {
		return []ProcessWindows{}, err
	}
	out := string(outBuffer.String())

	if errBuffer.Len() != 0 {
		err := errBuffer.String()
		return []ProcessWindows{}, errors.New(err)
	}

	processes, err := parseProcessesWindows(out)
	if err != nil {
		return []ProcessWindows{}, err
	}

	return processes, nil
}
