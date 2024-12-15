package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func VirtualGo(cfg config.State, align config.Alignment) []Segment {
	env, _ := os.LookupEnv("VIRTUALGO")
	if env == "" {
		return []Segment{}
	}

	return []Segment{{
		Name:       "vgo",
		Content:    env,
		Foreground: cfg.Theme.VirtualGoFg,
		Background: cfg.Theme.VirtualGoBg,
	}}
}
