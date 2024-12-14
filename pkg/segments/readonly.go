//go:build !windows
// +build !windows

package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/unix"
)

func Perms(cfg config.Config, align config.Alignment) []Segment {
	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("could not determine current working directory")
	}

	if unix.Access(cwd, unix.W_OK) == nil {
		return []Segment{}
	}
	return []Segment{{
		Name:       "perms",
		Content:    cfg.Symbols().Lock,
		Foreground: cfg.SelectedTheme().ReadonlyFg,
		Background: cfg.SelectedTheme().ReadonlyBg,
	}}
}
