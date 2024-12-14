package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Newline(cfg config.Config, align config.Alignment) []Segment {
	return []Segment{{NewLine: true}}
}
