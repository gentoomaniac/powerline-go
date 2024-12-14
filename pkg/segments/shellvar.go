package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func ShellVar(theme config.Theme) []segment {
	shellVarName := p.cfg.ShellVar
	varContent, varExists := os.LookupEnv(shellVarName)

	if !varExists {
		if shellVarName != "" {
			warn("Shell variable " + shellVarName + " does not exist.")
		}
		return []segment{}
	}

	if varContent == "" {
		if !p.cfg.ShellVarNoWarnEmpty {
			warn("Shell variable " + shellVarName + " is empty.")
		}
		return []segment{}
	}

	return []segment{{
		Name:       "shell-var",
		Content:    varContent,
		Foreground: theme.ShellVarFg,
		Background: theme.ShellVarBg,
	}}
}
