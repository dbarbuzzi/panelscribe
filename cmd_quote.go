package main

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func buildQuoteGroups(groupCounts string, border string) (string, error) {
	var groups []string

	if border == "top" || border == "both" {
		groups = append(groups, "---")
	}

	for i, gc := range groupCounts {
		if i > 0 {
			groups = append(groups, "---")
		}
		count, err := strconv.Atoi(string(gc))
		if err != nil {
			return "", fmt.Errorf("failed to convert '%s' to int: %w", string(gc), err)
		}
		for i := 0; i < count; i++ {
			groups = append(groups, "> ")
		}
	}

	if border == "bottom" || border == "both" {
		groups = append(groups, "---")
	}

	res := strings.Join(groups, "\n\n")
	return res, nil
}

func pipeToCommand(data []byte, cmd string) error {
	copyCmd := exec.Command(cmd)

	inPipe, err := copyCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to retrieve stdin pipe: %w", err)
	}

	if err := copyCmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	if _, err := inPipe.Write(data); err != nil {
		return fmt.Errorf("failed to write to stdin pipe: %w", err)
	}

	if err := inPipe.Close(); err != nil {
		return fmt.Errorf("failed to close stdin pipe: %w", err)
	}

	return copyCmd.Wait()
}

func executeCopy(data []byte) error {
	var cmd string
	switch runtime.GOOS {
	case "windows":
		cmd = "clip"
	case "darwin":
		cmd = "pbcopy"
	default:
		return errors.New("unsupported platform")
	}

	return pipeToCommand(data, cmd)
}

func validateQuoteGroupArg(arg string) error {
	if _, err := strconv.Atoi(arg); err != nil {
		return fmt.Errorf("arg must be an integer, got %s", arg)
	}

	for i, r := range arg {
		if r == '0' {
			return fmt.Errorf("arg should only contain positive integers, had '0' at position %d", i)
		}
	}

	return nil
}

func quoteAction(quoteGroups string, border string, print bool) error {
	if err := validateQuoteGroupArg(quoteGroups); err != nil {
		return fmt.Errorf("invalid args: %w", err)
	}

	qg, err := buildQuoteGroups(quoteGroups, border)
	if err != nil {
		return fmt.Errorf("failed to build quote groups: %w", err)
	}

	if print {
		fmt.Println(qg)
	} else {
		err := executeCopy([]byte(qg))
		if err != nil {
			return fmt.Errorf("failed to copy: %w", err)
		}
	}

	return nil
}
