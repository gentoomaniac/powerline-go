package segments

import (
	"os"
	"path"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Perlbrew(cfg config.Config) []Segment {
	env, _ := os.LookupEnv("PERLBREW_PERL")
	if env == "" {
		return []Segment{}
	}

	envName := path.Base(env)
	return []Segment{{
		Name:       "perlbrew",
		Content:    envName,
		Foreground: cfg.SelectedTheme().PerlbrewFg,
		Background: cfg.SelectedTheme().PerlbrewBg,
	}}
}
