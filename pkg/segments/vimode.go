//go:build broken

package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

func ViMode(theme config.Theme) []Segment {
	mode := p.cfg.ViMode
	if mode == "" {
		log.Warn().Msg("'--vi-mode' is not set.")
		return []Segment{}
	}

	switch mode {
	case "vicmd":
		return []Segment{{
			Name:       "vi-mode",
			Content:    "C",
			Foreground: theme.ViModeCommandFg,
			Background: theme.ViModeCommandBg,
		}}
	default: // usually "viins" or "main"
		return []Segment{{
			Name:       "vi-mode",
			Content:    "I",
			Foreground: theme.ViModeInsertFg,
			Background: theme.ViModeInsertBg,
		}}
	}
}
