package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func User(cfg config.Config, align config.Alignment) []Segment {
	var userPrompt string
	switch cfg.Shell {
	case "bash":
		userPrompt = "\\u"
	case "zsh":
		userPrompt = "%n"
	default:
		userPrompt = os.Getenv("USER")
	}

	var background uint8
	if userIsAdmin() {
		background = cfg.SelectedTheme().UsernameRootBg
	} else {
		background = cfg.SelectedTheme().UsernameBg
	}

	return []Segment{{
		Name:       "user",
		Content:    userPrompt,
		Foreground: cfg.SelectedTheme().UsernameFg,
		Background: background,
	}}
}
