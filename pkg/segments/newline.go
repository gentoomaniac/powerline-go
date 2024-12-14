package segments

import (
	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func Newline(theme config.Theme) []segment {
	return []segment{{NewLine: true}}
}
