package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func User(conf config.Config) []Segment {
	var userPrompt string
	switch conf.Shell {
	case "bash":
		userPrompt = "\\u"
	case "zsh":
		userPrompt = "%n"
	default:
		userPrompt = os.Getenv("USER")
	}

	var background uint8
	if userIsAdmin() {
		background = conf.SelectedTheme().UsernameRootBg
	} else {
		background = conf.SelectedTheme().UsernameBg
	}

	return []Segment{{
		Name:       "user",
		Content:    userPrompt,
		Foreground: conf.SelectedTheme().UsernameFg,
		Background: background,
	}}
}
