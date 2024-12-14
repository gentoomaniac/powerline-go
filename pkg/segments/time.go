//go:build broken

package segments

import (
	"strings"
	"time"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Time(theme config.Theme) []Segment {
	return []Segment{{
		Name:       "time",
		Content:    time.Now().Format(strings.TrimSpace(p.cfg.Time)),
		Foreground: theme.TimeFg,
		Background: theme.TimeBg,
	}}
}
