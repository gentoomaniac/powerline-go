package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func VirtualGo(theme config.Theme) []segment {
	env, _ := os.LookupEnv("VIRTUALGO")
	if env == "" {
		return []segment{}
	}

	return []segment{{
		Name:       "vgo",
		Content:    env,
		Foreground: theme.VirtualGoFg,
		Background: theme.VirtualGoBg,
	}}
}
