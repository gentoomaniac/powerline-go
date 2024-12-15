package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

func ViMode(cfg config.Config, align config.Alignment) []Segment {
	mode := cfg.ViMode
	if mode == "" {
		log.Warn().Msg("'--vi-mode' is not set.")
		return []Segment{}
	}

	switch mode {
	case "vicmd":
		return []Segment{{
			Name:       "vi-mode",
			Content:    "C",
			Foreground: cfg.Theme.ViModeCommandFg,
			Background: cfg.Theme.ViModeCommandBg,
		}}
	default: // usually "viins" or "main"
		return []Segment{{
			Name:       "vi-mode",
			Content:    "I",
			Foreground: cfg.Theme.ViModeInsertFg,
			Background: cfg.Theme.ViModeInsertBg,
		}}
	}
}
