package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func DotEnv(theme config.Theme) []segment {
	files := []string{".env", ".envrc"}
	dotEnv := false
	for _, file := range files {
		stat, err := os.Stat(file)
		if err == nil && !stat.IsDir() {
			dotEnv = true
			break
		}
	}
	if !dotEnv {
		return []segment{}
	}
	return []segment{{
		Name:       "dotenv",
		Content:    "\u2235",
		Foreground: theme.DotEnvFg,
		Background: theme.DotEnvBg,
	}}
}
