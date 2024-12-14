package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func Root(theme config.Theme) []segment {
	var foreground, background uint8
	if p.cfg.PrevError == 0 || p.cfg.StaticPromptIndicator {
		foreground = theme.CmdPassedFg
		background = theme.CmdPassedBg
	} else {
		foreground = theme.CmdFailedFg
		background = theme.CmdFailedBg
	}

	return []segment{{
		Name:       "root",
		Content:    p.shell.RootIndicator,
		Foreground: foreground,
		Background: background,
	}}
}
