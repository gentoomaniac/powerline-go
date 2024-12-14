package segments

import (
	"strings"
	"time"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func Time(theme config.Theme) []segment {
	return []segment{{
		Name:       "time",
		Content:    time.Now().Format(strings.TrimSpace(p.cfg.Time)),
		Foreground: theme.TimeFg,
		Background: theme.TimeBg,
	}}
}
