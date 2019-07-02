package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func buildQuoteGroups(groupCounts string, flip bool) (string, error) {
	var groups []string

	for _, gc := range groupCounts {
		groups = append(groups, "---")
		count, err := strconv.Atoi(string(gc))
		if err != nil {
			return "", err
		}
		for i := 0; i < count; i++ {
			groups = append(groups, "> ")
		}
	}

	if flip {
		// cut off first entry of "---" and “move” it to the end
		groups = append(groups[1:], "---")
	}

	res := fmt.Sprintf("\n%s\n", strings.Join(groups, "\n\n"))
	return res, nil
}

func executePbcopy(data []byte) error {
	copyCmd := exec.Command("pbcopy")

	inPipe, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}

	if _, err := inPipe.Write(data); err != nil {
		return err
	}

	if err := inPipe.Close(); err != nil {
		return err
	}

	return copyCmd.Wait()
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

func quoteAction(quoteGroups string, flip, print bool) error {
	if err := validateQuoteGroupArg(quoteGroups); err != nil {
		return err
	}

	qg, err := buildQuoteGroups(quoteGroups, flip)
	if err != nil {
		return err
	}

	if print {
		fmt.Println(qg)
	} else {
		err := executePbcopy([]byte(qg))
		if err != nil {
			return err
		}
	}

	return nil
}
