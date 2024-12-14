//go:build windows
// +build windows

package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Perms(cfg config.Config, align config.Alignment) []Segment {
	cwd := getValidCwd()
	const W_USR = 0o002
	// Check user's permissions on directory in a portable but probably slower way
	fileInfo, _ := os.Stat(cwd)
	if fileInfo.Mode()&W_USR == W_USR {
		return []Segment{}
	}
	return []Segment{{
		Name:       "perms",
		Content:    config.Defaults.Symbols().Lock,
		Foreground: config.Defaults.SelectedTheme().ReadonlyFg,
		Background: config.Defaults.SelectedTheme().ReadonlyBg,
	}}
}
