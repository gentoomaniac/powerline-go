//go:build broken

package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Root(theme config.Theme) []Segment {
	var foreground, background uint8
	if p.cfg.PrevError == 0 || p.cfg.StaticPromptIndicator {
		foreground = theme.CmdPassedFg
		background = theme.CmdPassedBg
	} else {
		foreground = theme.CmdFailedFg
		background = theme.CmdFailedBg
	}

	return []Segment{{
		Name:       "root",
		Content:    p.shell.RootIndicator,
		Foreground: foreground,
		Background: background,
	}}
}
