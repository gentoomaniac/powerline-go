//go:build windows
// +build windows

package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Perms(theme config.Theme) []Segment {
	cwd := p.cwd
	const W_USR = 0002
	// Check user's permissions on directory in a portable but probably slower way
	fileInfo, _ := os.Stat(cwd)
	if fileInfo.Mode()&W_USR == W_USR {
		return []Segment{}
	}
	return []Segment{{
		Name:       "perms",
		Content:    p.symbols.Lock,
		Foreground: theme.ReadonlyFg,
		Background: theme.ReadonlyBg,
	}}
}
