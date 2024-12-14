//go:build broken

package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

func ShellVar(cfg config.Config, align config.Alignment) []Segment {
	shellVarName := p.cfg.ShellVar
	varContent, varExists := os.LookupEnv(shellVarName)

	if !varExists {
		if shellVarName != "" {
			log.Warn().Msgf("Shell variable %s does not exist.", shellVarName)
		}
		return []Segment{}
	}

	if varContent == "" {
		if !p.cfg.ShellVarNoWarnEmpty {
			log.Warn().Msgf("Shell variable %s is empty.", shellVarName)
		}
		return []Segment{}
	}

	return []Segment{{
		Name:       "shell-var",
		Content:    varContent,
		Foreground: cfg.SelectedTheme().ShellVarFg,
		Background: cfg.SelectedTheme().ShellVarBg,
	}}
}
