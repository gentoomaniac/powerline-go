package segments

import (
	"fmt"
	"runtime"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	"github.com/shirou/gopsutil/v3/load"
)

func Load(cfg config.State, align config.Alignment) []Segment {
	c := runtime.NumCPU()
	a, err := load.Avg()
	if err != nil {
		return []Segment{}
	}
	bg := cfg.Theme.LoadBg

	load := a.Load5
	switch cfg.Theme.LoadAvgValue {
	case 1:
		load = a.Load1
	case 15:
		load = a.Load15
	}

	if load > float64(c)*cfg.Theme.LoadThresholdBad {
		bg = cfg.Theme.LoadHighBg
	}

	return []Segment{{
		Name:       "load",
		Content:    fmt.Sprintf("%.2f", a.Load5),
		Foreground: cfg.Theme.LoadFg,
		Background: bg,
	}}
}
