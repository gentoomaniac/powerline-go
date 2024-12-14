package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func AWS(cfg config.Config) []Segment {
	profile := os.Getenv("AWS_PROFILE")
	region := os.Getenv("AWS_DEFAULT_REGION")
	if profile == "" {
		profile = os.Getenv("AWS_VAULT")
		if profile == "" {
			return []Segment{}
		}
	}
	var r string
	if region != "" {
		r = " (" + region + ")"
	}
	return []Segment{{
		Name:       "aws",
		Content:    profile + r,
		Foreground: cfg.SelectedTheme().AWSFg,
		Background: cfg.SelectedTheme().AWSBg,
	}}
}
