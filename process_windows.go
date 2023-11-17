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
