package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func PlEnv(theme config.Theme) []Segment {
	env, _ := os.LookupEnv("PLENV_VERSION")
	if env == "" {
		return []Segment{}
	}
	return []Segment{{
		Name:       "plenv",
		Content:    env,
		Foreground: theme.PlEnvFg,
		Background: theme.PlEnvBg,
	}}
}
