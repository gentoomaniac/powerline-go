package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func User(cfg config.State, align config.Alignment) []Segment {
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
		background = cfg.Theme.UsernameRootBg
	} else {
		background = cfg.Theme.UsernameBg
	}

	return []Segment{{
		Name:       "user",
		Content:    userPrompt,
		Foreground: cfg.Theme.UsernameFg,
		Background: background,
	}}
}
