//go:build windows
// +build windows

package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func Perms(theme config.Theme) []segment {
	cwd := p.cwd
	const W_USR = 0002
	// Check user's permissions on directory in a portable but probably slower way
	fileInfo, _ := os.Stat(cwd)
	if fileInfo.Mode()&W_USR == W_USR {
		return []segment{}
	}
	return []segment{{
		Name:       "perms",
		Content:    p.symbols.Lock,
		Foreground: theme.ReadonlyFg,
		Background: theme.ReadonlyBg,
	}}
}
