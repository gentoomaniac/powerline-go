package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Root(cfg config.Config, align config.Alignment) []Segment {
	var foreground, background uint8
	if cfg.PrevError == 0 || cfg.StaticPromptIndicator {
		foreground = cfg.SelectedTheme().CmdPassedFg
		background = cfg.SelectedTheme().CmdPassedBg
	} else {
		foreground = cfg.SelectedTheme().CmdFailedFg
		background = cfg.SelectedTheme().CmdFailedBg
	}

	return []Segment{{
		Name:       "root",
		Content:    cfg.CurrentShell().RootIndicator,
		Foreground: foreground,
		Background: background,
	}}
}
