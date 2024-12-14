package segments

import (
	"fmt"
	"runtime"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	"github.com/shirou/gopsutil/v3/load"
)

func Load(cfg config.Config) []Segment {
	c := runtime.NumCPU()
	a, err := load.Avg()
	if err != nil {
		return []Segment{}
	}
	bg := cfg.SelectedTheme().LoadBg

	load := a.Load5
	switch cfg.SelectedTheme().LoadAvgValue {
	case 1:
		load = a.Load1
	case 15:
		load = a.Load15
	}

	if load > float64(c)*cfg.SelectedTheme().LoadThresholdBad {
		bg = cfg.SelectedTheme().LoadHighBg
	}

	return []Segment{{
		Name:       "load",
		Content:    fmt.Sprintf("%.2f", a.Load5),
		Foreground: cfg.SelectedTheme().LoadFg,
		Background: bg,
	}}
}
