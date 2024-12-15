//go:build !windows
// +build !windows

package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
	"golang.org/x/sys/unix"
)

func Perms(cfg config.Config, align config.Alignment) []Segment {
	if unix.Access(cfg.Cwd, unix.W_OK) == nil {
		return []Segment{}
	}
	return []Segment{{
		Name:       "perms",
		Content:    cfg.Symbols().Lock,
		Foreground: cfg.Theme.ReadonlyFg,
		Background: cfg.Theme.ReadonlyBg,
	}}
}
