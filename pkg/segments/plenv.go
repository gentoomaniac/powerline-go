package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func PlEnv(theme config.Theme) []segment {
	env, _ := os.LookupEnv("PLENV_VERSION")
	if env == "" {
		return []segment{}
	}
	return []segment{{
		Name:       "plenv",
		Content:    env,
		Foreground: theme.PlEnvFg,
		Background: theme.PlEnvBg,
	}}
}
