package segments

import (
	"fmt"
	"runtime"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	"github.com/shirou/gopsutil/v3/load"
)

func Load(theme config.Theme) []Segment {
	c := runtime.NumCPU()
	a, err := load.Avg()
	if err != nil {
		return []Segment{}
	}
	bg := theme.LoadBg

	load := a.Load5
	switch theme.LoadAvgValue {
	case 1:
		load = a.Load1
	case 15:
		load = a.Load15
	}

	if load > float64(c)*theme.LoadThresholdBad {
		bg = theme.LoadHighBg
	}

	return []Segment{{
		Name:       "load",
		Content:    fmt.Sprintf("%.2f", a.Load5),
		Foreground: theme.LoadFg,
		Background: bg,
	}}
}
