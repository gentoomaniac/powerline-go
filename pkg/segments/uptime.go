package segments

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

func Uptime(cfg config.State, align config.Alignment) []Segment {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		log.Error().Err(err).Msg("failed getting uptime")
	}

	seconds, err := strconv.ParseFloat(strings.TrimSpace(strings.Split(string(data), " ")[0]), 64)
	if err != nil {
		log.Error().Err(err).Msg("failed parsing uptime")
	}
	return []Segment{{
		Name:       "time",
		Content:    time.Duration(seconds * float64(time.Second)).String(),
		Foreground: cfg.Theme.TimeFg,
		Background: cfg.Theme.TimeBg,
	}}
}
