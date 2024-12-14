//go:build broken

package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

func ViMode(cfg config.Config, align config.Alignment) []Segment {
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
			Foreground: cfg.SelectedTheme().ViModeCommandFg,
			Background: cfg.SelectedTheme().ViModeCommandBg,
		}}
	default: // usually "viins" or "main"
		return []Segment{{
			Name:       "vi-mode",
			Content:    "I",
			Foreground: cfg.SelectedTheme().ViModeInsertFg,
			Background: cfg.SelectedTheme().ViModeInsertBg,
		}}
	}
}
