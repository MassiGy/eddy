package main

import (
	"bytes"
	"os/exec"
)

const target_shell = "bash"

func run_command(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(target_shell, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
func eval(command string) (buffer [][]rune) {
	buffer = [][]rune{}
	stdout, stderr, err := run_command(command)

	buffer = append(buffer, []rune{})
	for _, c := range command {
		buffer[0] = append(buffer[0], rune(c))
	}

	buffer = append(buffer, []rune{})
	buffer = append(buffer, []rune{})
	buffer = append(buffer, []rune{'O', 'u', 't', 'p', 'u', 't', ':'})
	buffer = append(buffer, []rune{'=', '=', '=', '=', '=', '=', '='})
	buffer = append(buffer, []rune{})

	buffer = append(buffer, []rune{})
	buffer = append(buffer, []rune{'S', 't', 'd', 'o', 'u', 't', ':'})
	buffer = append(buffer, []rune{})
	l := len(buffer)
	for _, c := range stdout {
		if c == '\n' {
			buffer = append(buffer, []rune{})
			l++
		}
		buffer[l-1] = append(buffer[l-1], rune(c))
	}

	buffer = append(buffer, []rune{})
	buffer = append(buffer, []rune{'S', 't', 'd', 'e', 'r', 'r', ':'})
	buffer = append(buffer, []rune{})
	l = len(buffer)
	for _, c := range stderr {
		if c == '\n' {
			buffer = append(buffer, []rune{})
			l++
		}
		buffer[l-1] = append(buffer[l-1], rune(c))
	}

	buffer = append(buffer, []rune{})
	buffer = append(buffer, []rune{'E', 'r', 'r', ':'})
	buffer = append(buffer, []rune{})
	l = len(buffer)
	if err != nil {
		for _, c := range err.Error() {
			if c == '\n' {
				buffer = append(buffer, []rune{})
				l++
			}
			buffer[l-1] = append(buffer[l-1], rune(c))
		}
	}

	return buffer
}
