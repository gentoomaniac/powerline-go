package segments

import (
	"strings"
	"time"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Time(cfg config.Config, align config.Alignment) []Segment {
	return []Segment{{
		Name: "time",
		// TODO: use go to get the time and instead use config to provide a format string
		Content:    time.Now().Format(strings.TrimSpace(cfg.Time)),
		Foreground: cfg.SelectedTheme().TimeFg,
		Background: cfg.SelectedTheme().TimeBg,
	}}
}
