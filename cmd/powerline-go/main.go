package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/gentoomaniac/powerline-go/pkg/logging"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
	"github.com/rs/zerolog/log"
)

var (
	version = "unknown"
	commit  = "unknown"
	binName = "unknown"
	builtBy = "unknown"
	date    = "unknown"
)

var cli config.Config

func main() {
	ctx := kong.Parse(&cli,
		kong.UsageOnError(),
		kong.Configuration(kong.JSON, filepath.Join(config.ConfigPath(), config.ConfigFileName)),
		kong.Vars{
			"version": version,
			"commit":  commit,
			"binName": binName,
			"builtBy": builtBy,
			"date":    date,
		},
	)

	cli.Json = false
	logging.Setup(&cli.LoggingConfig)

	state := config.NewStateFromConfig(cli)

	if cli.SaveConfig {
		if err := cli.Save(); err != nil {
			log.Fatal().Err(err).Msg("could not write default config")
		}
		ctx.Exit(0)
	}

	var exists bool
	state.Theme, exists = config.DefaultThemes[cli.Theme]
	if !exists {
		theme, err := config.ThemeFromFile(cli.Theme)
		if err != nil {
			log.Warn().Err(err).Msg("could't load theme, falling back to default")
			state.Theme = config.DefaultThemes["default"]
		} else {
			state.Theme = *theme
		}
	}

	if strings.HasSuffix(state.Mode, ".json") {
		file, err := os.ReadFile(state.Mode)
		if err == nil {
			symbols := state.Modes["compatible"]
			err = json.Unmarshal(file, &symbols)
			if err == nil {
				state.Modes[state.Mode] = symbols
			} else {
				log.Error().Err(err).Msg("error reading mode")
			}
		}
	}

	p := pwl.NewPowerline(state, config.AlignLeft)
	if p.SupportsRightModules() && p.HasRightModules() && !state.Eval {
		log.Fatal().Msg("'--modules-right' requires '--eval' mode")
	}

	fmt.Print(p.Render())

	ctx.Exit(0)
}
