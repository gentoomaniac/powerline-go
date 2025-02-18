package segments

import (
	"os"
	"path"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Perlbrew(cfg config.State, align config.Alignment) []Segment {
	env, _ := os.LookupEnv("PERLBREW_PERL")
	if env == "" {
		return []Segment{}
	}

	envName := path.Base(env)
	return []Segment{{
		Name:       "perlbrew",
		Content:    envName,
		Foreground: cfg.Theme.PerlbrewFg,
		Background: cfg.Theme.PerlbrewBg,
	}}
}
