package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Root(cfg config.State, align config.Alignment) []Segment {
	var foreground, background uint8
	if cfg.PrevError == 0 || cfg.StaticPromptIndicator {
		foreground = cfg.Theme.CmdPassedFg
		background = cfg.Theme.CmdPassedBg
	} else {
		foreground = cfg.Theme.CmdFailedFg
		background = cfg.Theme.CmdFailedBg
	}

	return []Segment{{
		Name:       "root",
		Content:    cfg.CurrentShell().RootIndicator,
		Foreground: foreground,
		Background: background,
	}}
}
