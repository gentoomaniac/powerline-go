package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func ViMode(theme config.Theme) []segment {
	mode := p.cfg.ViMode
	if mode == "" {
		warn("'--vi-mode' is not set.")
		return []segment{}
	}

	switch mode {
	case "vicmd":
		return []segment{{
			Name:       "vi-mode",
			Content:    "C",
			Foreground: theme.ViModeCommandFg,
			Background: theme.ViModeCommandBg,
		}}
	default: // usually "viins" or "main"
		return []segment{{
			Name:       "vi-mode",
			Content:    "I",
			Foreground: theme.ViModeInsertFg,
			Background: theme.ViModeInsertBg,
		}}
	}
}
