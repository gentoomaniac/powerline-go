package segments

import (
	"strconv"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Jobs(cfg config.State, align config.Alignment) []Segment {
	if cfg.Jobs <= 0 {
		return []Segment{}
	}
	return []Segment{{
		Name:       "jobs",
		Content:    strconv.Itoa(cfg.Jobs),
		Foreground: cfg.Theme.JobsFg,
		Background: cfg.Theme.JobsBg,
	}}
}
