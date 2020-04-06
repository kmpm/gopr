package cmd

import (
	"runtime"

	"github.com/kmpm/goflip/lib/shell"
)

func getShell(userShell string) (string, error) {
	if userShell != "" {
		return userShell, nil
	}
	return shell.Detect()
}

func runtimeOS() string {
	return runtime.GOOS
}
