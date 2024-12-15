package segments

import (
	"strings"
	"time"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Time(cfg config.Config, align config.Alignment) []Segment {
	return []Segment{{
		Name:       "time",
		Content:    time.Now().Format(strings.TrimSpace(cfg.Time)),
		Foreground: cfg.Theme.TimeFg,
		Background: cfg.Theme.TimeBg,
	}}
}
