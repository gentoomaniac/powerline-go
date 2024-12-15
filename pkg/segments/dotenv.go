package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func DotEnv(cfg config.Config, align config.Alignment) []Segment {
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
		return []Segment{}
	}
	return []Segment{{
		Name:       "dotenv",
		Content:    "\u2235",
		Foreground: cfg.Theme.DotEnvFg,
		Background: cfg.Theme.DotEnvBg,
	}}
}
