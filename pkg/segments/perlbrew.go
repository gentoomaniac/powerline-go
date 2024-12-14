package segments

import (
	"os"
	"path"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Perlbrew(theme config.Theme) []Segment {
	env, _ := os.LookupEnv("PERLBREW_PERL")
	if env == "" {
		return []Segment{}
	}

	envName := path.Base(env)
	return []Segment{{
		Name:       "perlbrew",
		Content:    envName,
		Foreground: theme.PerlbrewFg,
		Background: theme.PerlbrewBg,
	}}
}
