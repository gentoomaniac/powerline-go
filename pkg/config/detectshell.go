package config

import (
	"path"
	"strings"
)

func detectShell(shellExe string) string {
	var shell string
	shellExe = path.Base(shellExe)
	if strings.Contains(shellExe, "bash") {
		shell = "bash"
	} else if strings.Contains(shellExe, "zsh") {
		shell = "zsh"
	} else if strings.Contains(shellExe, "go") { // for testing with `go run`
		shell = "bash"
	} else {
		shell = "bare"
	}
	return shell
}
