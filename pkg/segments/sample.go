package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
)

// Register your segment here: https://github.com/gentoomaniac/powerline-go/blob/main/pkg/powerline/modules.go

func Sample(cfg config.State, align config.Alignment) []Segment {
	return []Segment{{
		// currently unused?
		Name: "sample",
		// The text that will be rendered in the segment
		Content: "fizzbuzz",
		// `config.Theme` currently needs to be extended for segment specific colours
		// They should be added to the Default themes as well.
		Foreground: cfg.Theme.AWSFg,
		Background: cfg.Theme.AWSBg,
	}}
}
