package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/gentoomaniac/powerline-go/pkg/logging"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
	"github.com/rs/zerolog/log"
)

var cfg config.Config

func main() {
	ctx := kong.Parse(&cfg, kong.UsageOnError(), kong.Configuration(kong.JSON, config.ConfigPath()))

	if cfg.SaveDefaultConfig {
		if err := cfg.EnsureExists(); err != nil {
			log.Fatal().Msg("could not write default config")
		}
		ctx.Exit(0)
	}

	cfg.Json = false
	logging.Setup(&cfg.LoggingConfig)

	if strings.HasSuffix(cfg.Theme, ".json") {
		file, err := ioutil.ReadFile(cfg.Theme)
		if err == nil {
			theme := cfg.Themes[config.Defaults.Theme]
			err = json.Unmarshal(file, &theme)
			if err == nil {
				cfg.Themes[cfg.Theme] = theme
			} else {
				log.Error().Err(err).Msg("error reading theme")
			}
		}
	}

	if strings.HasSuffix(cfg.Mode, ".json") {
		file, err := ioutil.ReadFile(cfg.Mode)
		if err == nil {
			symbols := cfg.Modes[config.Defaults.Mode]
			err = json.Unmarshal(file, &symbols)
			if err == nil {
				cfg.Modes[cfg.Mode] = symbols
			} else {
				log.Error().Err(err).Msg("error reading mode")
			}
		}
	}

	p := pwl.NewPowerline(cfg, pwl.AlignLeft)
	if p.SupportsRightModules() && p.HasRightModules() && !cfg.Eval {
		log.Fatal().Msg("'--modules-right' requires '--eval' mode")
	}

	fmt.Print(p.Render())

	ctx.Exit(0)
}
