package segments

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

const (
	micro  rune = '\u00B5'
	milli  rune = 'm'
	second rune = 's'
	minute rune = 'm'
	hour   rune = 'h'
)

const (
	nanoseconds  int64 = 1
	microseconds int64 = nanoseconds * 1000
	milliseconds int64 = microseconds * 1000
	seconds      int64 = milliseconds * 1000
	minutes      int64 = seconds * 60
	hours        int64 = minutes * 60
)

func Duration(cfg config.Config, align config.Alignment) []Segment {
	if cfg.Duration == "" {
		return []Segment{{
			Name:       "duration",
			Content:    "No duration",
			Foreground: cfg.Theme.DurationFg,
			Background: cfg.Theme.DurationBg,
		}}
	}

	durationValue := strings.Trim(cfg.Duration, "'\"")
	durationMinValue := strings.Trim(cfg.DurationMin, "'\"")

	hasPrecision := strings.Contains(durationValue, ".")

	durationFloat, err := strconv.ParseFloat(durationValue, 64)
	durationMinFloat, _ := strconv.ParseFloat(durationMinValue, 64)
	if err != nil {
		return []Segment{{
			Name:       "duration",
			Content:    fmt.Sprintf("Failed to convert '%s' to a number", cfg.Duration),
			Foreground: cfg.Theme.DurationFg,
			Background: cfg.Theme.DurationBg,
		}}
	}

	if durationFloat < durationMinFloat {
		return []Segment{}
	}

	duration := time.Duration(durationFloat * float64(time.Second.Nanoseconds()))

	if duration <= 0 {
		return []Segment{}
	}

	var content string
	ns := duration.Nanoseconds()
	if ns > hours {
		hrs := ns / hours
		ns -= hrs * hours
		mins := ns / minutes
		content = fmt.Sprintf("%d%c %d%c", hrs, hour, mins, minute)
	} else if ns > minutes {
		mins := ns / minutes
		ns -= mins * minutes
		secs := ns / seconds
		content = fmt.Sprintf("%d%c %d%c", mins, minute, secs, second)
	} else if !hasPrecision {
		secs := ns / seconds
		content = fmt.Sprintf("%d%c", secs, second)
	} else if ns > seconds {
		secs := ns / seconds
		ns -= secs * seconds
		millis := ns / milliseconds
		content = fmt.Sprintf("%d%c %d%c%c", secs, second, millis, milli, second)
	} else if ns > milliseconds || cfg.DurationLowPrecision {
		millis := ns / milliseconds
		ns -= millis * milliseconds
		micros := ns / microseconds
		if cfg.DurationLowPrecision {
			content = fmt.Sprintf("%d%c%c", millis, milli, second)
		} else {
			content = fmt.Sprintf("%d%c%c %d%c%c", millis, milli, second, micros, micro, second)
		}
	} else {
		content = fmt.Sprintf("%d%c%c", ns/microseconds, micro, second)
	}

	return []Segment{{
		Name:       "duration",
		Content:    content,
		Foreground: cfg.Theme.DurationFg,
		Background: cfg.Theme.DurationBg,
	}}
}
