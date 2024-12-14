package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func ShEnv(theme config.Theme) []Segment {
	env, _ := os.LookupEnv("SHENV_VERSION")
	if env == "" {
		return []Segment{}
	}
	return []Segment{{
		Name:       "shenv",
		Content:    env,
		Foreground: theme.ShEnvFg,
		Background: theme.ShEnvBg,
	}}
}
