package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func ShEnv(cfg config.State, align config.Alignment) []Segment {
	env, _ := os.LookupEnv("SHENV_VERSION")
	if env == "" {
		return []Segment{}
	}
	return []Segment{{
		Name:       "shenv",
		Content:    env,
		Foreground: cfg.Theme.ShEnvFg,
		Background: cfg.Theme.ShEnvBg,
	}}
}
