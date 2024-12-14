package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func VirtualGo(theme config.Theme) []Segment {
	env, _ := os.LookupEnv("VIRTUALGO")
	if env == "" {
		return []Segment{}
	}

	return []Segment{{
		Name:       "vgo",
		Content:    env,
		Foreground: theme.VirtualGoFg,
		Background: theme.VirtualGoBg,
	}}
}
