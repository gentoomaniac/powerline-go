package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func ShEnv(cfg config.Config, align config.Alignment) []Segment {
	env, _ := os.LookupEnv("SHENV_VERSION")
	if env == "" {
		return []Segment{}
	}
	return []Segment{{
		Name:       "shenv",
		Content:    env,
		Foreground: cfg.SelectedTheme().ShEnvFg,
		Background: cfg.SelectedTheme().ShEnvBg,
	}}
}
