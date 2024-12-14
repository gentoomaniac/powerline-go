//go:build !windows
// +build !windows

package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
	"golang.org/x/sys/unix"
)

func Perms(theme config.Theme) []segment {
	cwd := p.cwd
	if unix.Access(cwd, unix.W_OK) == nil {
		return []segment{}
	}
	return []segment{{
		Name:       "perms",
		Content:    p.symbols.Lock,
		Foreground: theme.ReadonlyFg,
		Background: theme.ReadonlyBg,
	}}
}
