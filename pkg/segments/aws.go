package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func AWS(theme config.Theme) []segment {
	profile := os.Getenv("AWS_PROFILE")
	region := os.Getenv("AWS_DEFAULT_REGION")
	if profile == "" {
		profile = os.Getenv("AWS_VAULT")
		if profile == "" {
			return []segment{}
		}
	}
	var r string
	if region != "" {
		r = " (" + region + ")"
	}
	return []segment{{
		Name:       "aws",
		Content:    profile + r,
		Foreground: theme.AWSFg,
		Background: theme.AWSBg,
	}}
}
