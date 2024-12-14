//go:build broken

package segments

import (
	"strconv"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Jobs(theme config.Theme) []Segment {
	if p.cfg.Jobs <= 0 {
		return []Segment{}
	}
	return []Segment{{
		Name:       "jobs",
		Content:    strconv.Itoa(p.cfg.Jobs),
		Foreground: theme.JobsFg,
		Background: theme.JobsBg,
	}}
}
