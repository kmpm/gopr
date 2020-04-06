// +build !windows

package shell

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrUnknownShell = errors.New("unknown shell")
)

func detectShell() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		fmt.Printf("The default lines below are for a sh/bash shell, you can specify the shell you're using, with the --shell flag.\n\n")
		return "", ErrUnknownShell
	}
	return filepath.Base(shell), nil
}
