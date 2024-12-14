//go:build broken

package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func User(theme config.Theme) []Segment {
	var userPrompt string
	switch p.cfg.Shell {
	case "bash":
		userPrompt = "\\u"
	case "zsh":
		userPrompt = "%n"
	default:
		userPrompt = p.username
	}

	var background uint8
	if p.userIsAdmin {
		background = theme.UsernameRootBg
	} else {
		background = theme.UsernameBg
	}

	return []Segment{{
		Name:       "user",
		Content:    userPrompt,
		Foreground: theme.UsernameFg,
		Background: background,
	}}
}
