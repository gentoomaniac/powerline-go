package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

func comments(lines ...string) string {
	return " " + strings.Join(lines, "\n"+" ")
}

func commentsWithDefaults(lines ...string) string {
	return comments(lines...) + "\n"
}

func main() {
	logging.Setup(&logging.LoggingConfig{Verbosity: 4})
	flag.Parse()

	cfg := config.Defaults
	err := cfg.Load()
	if err != nil {
		log.Error().Err(err).Msg("error loading config")
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "cwd-mode":
			cfg.CwdMode = *args.CwdMode
		case "cwd-max-depth":
			cfg.CwdMaxDepth = *args.CwdMaxDepth
		case "cwd-max-dir-size":
			cfg.CwdMaxDirSize = *args.CwdMaxDirSize
		case "colorize-hostname":
			cfg.ColorizeHostname = *args.ColorizeHostname
		case "hostname-only-if-ssh":
			cfg.HostnameOnlyIfSSH = *args.HostnameOnlyIfSSH
		case "alternate-ssh-icon":
			cfg.SshAlternateIcon = *args.SshAlternateIcon
		case "east-asian-width":
			cfg.EastAsianWidth = *args.EastAsianWidth
		case "newline":
			cfg.PromptOnNewLine = *args.PromptOnNewLine
		case "static-prompt-indicator":
			cfg.StaticPromptIndicator = *args.StaticPromptIndicator
		case "venv-name-size-limit":
			cfg.VenvNameSizeLimit = *args.VenvNameSizeLimit
		case "jobs":
			cfg.Jobs = *args.Jobs
		case "git-assume-unchanged-size":
			cfg.GitAssumeUnchangedSize = *args.GitAssumeUnchangedSize
		case "git-disable-stats":
			cfg.GitDisableStats = strings.Split(*args.GitDisableStats, ",")
		case "git-mode":
			cfg.GitMode = *args.GitMode
		case "mode":
			cfg.Mode = *args.Mode
		case "theme":
			cfg.Theme = *args.Theme
		case "shell":
			cfg.Shell = *args.Shell
		case "modules":
			cfg.Modules = strings.Split(*args.Modules, ",")
		case "modules-right":
			cfg.ModulesRight = strings.Split(*args.ModulesRight, ",")
		case "priority":
			cfg.Priority = strings.Split(*args.Priority, ",")
		case "max-width":
			cfg.MaxWidthPercentage = *args.MaxWidthPercentage
		case "truncate-segment-width":
			cfg.TruncateSegmentWidth = *args.TruncateSegmentWidth
		case "error":
			cfg.PrevError = *args.PrevError
		case "numeric-exit-codes":
			cfg.NumericExitCodes = *args.NumericExitCodes
		case "ignore-repos":
			cfg.IgnoreRepos = strings.Split(*args.IgnoreRepos, ",")
		case "shorten-gke-names":
			cfg.ShortenGKENames = *args.ShortenGKENames
		case "shorten-eks-names":
			cfg.ShortenEKSNames = *args.ShortenEKSNames
		case "shorten-openshift-names":
			cfg.ShortenOpenshiftNames = *args.ShortenOpenshiftNames
		case "shell-var":
			cfg.ShellVar = *args.ShellVar
		case "shell-var-no-warn-empty":
			cfg.ShellVarNoWarnEmpty = *args.ShellVarNoWarnEmpty
		case "trim-ad-domain":
			cfg.TrimADDomain = *args.TrimADDomain
		case "path-aliases":
			for _, pair := range strings.Split(*args.PathAliases, ",") {
				kv := strings.SplitN(pair, "=", 2)
				cfg.PathAliases[kv[0]] = kv[1]
			}
		case "duration":
			cfg.Duration = *args.Duration
		case "duration-min":
			cfg.DurationMin = *args.DurationMin
		case "duration-low-precision":
			cfg.DurationLowPrecision = *args.DurationLowPrecision
		case "eval":
			cfg.Eval = *args.Eval
		case "condensed":
			cfg.Condensed = *args.Condensed
		case "ignore-warnings":
			cfg.IgnoreWarnings = *args.IgnoreWarnings
		case "time":
			cfg.Time = *args.Time
		case "vi-mode":
			cfg.ViMode = *args.ViMode
		}
	})

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

	p := pwl.NewPowerline(cfg, getValidCwd(), alignLeft)
	if p.SupportsRightModules() && p.HasRightModules() && !cfg.Eval {
		panic("Flag '-modules-right' requires '-eval' mode.")
	}

	fmt.Print(p.Draw())
}
