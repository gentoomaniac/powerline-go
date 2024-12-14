package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Newline(theme config.Theme) []Segment {
	return []Segment{{NewLine: true}}
}
