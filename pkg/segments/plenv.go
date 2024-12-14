package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func PlEnv(cfg config.Config) []Segment {
	env, _ := os.LookupEnv("PLENV_VERSION")
	if env == "" {
		return []Segment{}
	}
	return []Segment{{
		Name:       "plenv",
		Content:    env,
		Foreground: cfg.SelectedTheme().PlEnvFg,
		Background: cfg.SelectedTheme().PlEnvBg,
	}}
}
