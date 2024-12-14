//go:build broken

// //go:build !windows
// // +build !windows

package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
	"golang.org/x/sys/unix"
)

func Perms(theme config.Theme) []Segment {
	cwd := p.cwd
	if unix.Access(cwd, unix.W_OK) == nil {
		return []Segment{}
	}
	return []Segment{{
		Name:       "perms",
		Content:    p.symbols.Lock,
		Foreground: theme.ReadonlyFg,
		Background: theme.ReadonlyBg,
	}}
}
