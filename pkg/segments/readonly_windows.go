//go:build windows
// +build windows

package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Perms(cfg config.Config, align config.Alignment) []Segment {
	const W_USR = 0o002
	// Check user's permissions on directory in a portable but probably slower way
	fileInfo, _ := os.Stat(cfg.Cwd)
	if fileInfo.Mode()&W_USR == W_USR {
		return []Segment{}
	}
	return []Segment{{
		Name:       "perms",
		Content:    cfg.Symbols().Lock,
		Foreground: cfg.Theme.ReadonlyFg,
		Background: cfg.Theme.ReadonlyBg,
	}}
}
