package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func ShEnv(theme config.Theme) []segment {
	env, _ := os.LookupEnv("SHENV_VERSION")
	if env == "" {
		return []segment{}
	}
	return []segment{{
		Name:       "shenv",
		Content:    env,
		Foreground: theme.ShEnvFg,
		Background: theme.ShEnvBg,
	}}
}
