package segments

import (
	"strconv"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Jobs(cfg config.Config) []Segment {
	if cfg.Jobs <= 0 {
		return []Segment{}
	}
	return []Segment{{
		Name:       "jobs",
		Content:    strconv.Itoa(cfg.Jobs),
		Foreground: cfg.SelectedTheme().JobsFg,
		Background: cfg.SelectedTheme().JobsBg,
	}}
}
