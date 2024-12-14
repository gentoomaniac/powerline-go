package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/gentoomaniac/powerline-go/pkg/logging"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
	"github.com/rs/zerolog/log"
)

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func getValidCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		var exists bool
		cwd, exists = os.LookupEnv("PWD")
		if !exists {
			log.Warn().Msg("Your current directory is invalid.")
			print("> ")
			os.Exit(1)
		}
	}

	parts := strings.Split(cwd, string(os.PathSeparator))
	up := cwd

	for len(parts) > 0 && !pathExists(up) {
		parts = parts[:len(parts)-1]
		up = strings.Join(parts, string(os.PathSeparator))
	}
	if cwd != up {
		log.Warn().Msgf("Your current directory is invalid. Lowest valid directory: %s", up)
	}
	return cwd
}

var cfg = config.Defaults

func main() {
	ctx := kong.Parse(&cfg, kong.UsageOnError(), kong.Configuration(kong.JSON, config.ConfigPath()))

	logging.Setup(&logging.LoggingConfig{Verbosity: 4})

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

	p := pwl.NewPowerline(cfg, getValidCwd(), pwl.AlignLeft)
	if p.SupportsRightModules() && p.HasRightModules() && !cfg.Eval {
		panic("Flag '-modules-right' requires '-eval' mode.")
	}

	fmt.Print(p.Draw())

	ctx.Exit(0)
}
