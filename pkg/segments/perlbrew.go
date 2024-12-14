package segments

import (
	"os"
	"path"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func Perlbrew(theme config.Theme) []segment {
	env, _ := os.LookupEnv("PERLBREW_PERL")
	if env == "" {
		return []segment{}
	}

	envName := path.Base(env)
	return []segment{{
		Name:       "perlbrew",
		Content:    envName,
		Foreground: theme.PerlbrewFg,
		Background: theme.PerlbrewBg,
	}}
}
