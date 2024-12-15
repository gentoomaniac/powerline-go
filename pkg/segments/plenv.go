package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func PlEnv(cfg config.Config, align config.Alignment) []Segment {
	env, _ := os.LookupEnv("PLENV_VERSION")
	if env == "" {
		return []Segment{}
	}
	return []Segment{{
		Name:       "plenv",
		Content:    env,
		Foreground: cfg.Theme.PlEnvFg,
		Background: cfg.Theme.PlEnvBg,
	}}
}
