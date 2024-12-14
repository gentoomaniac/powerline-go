package segments

import (
	"strconv"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func Jobs(theme config.Theme) []segment {
	if p.cfg.Jobs <= 0 {
		return []segment{}
	}
	return []segment{{
		Name:       "jobs",
		Content:    strconv.Itoa(p.cfg.Jobs),
		Foreground: theme.JobsFg,
		Background: theme.JobsBg,
	}}
}
