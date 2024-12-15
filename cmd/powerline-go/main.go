package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/gentoomaniac/powerline-go/pkg/logging"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
	"github.com/rs/zerolog/log"
)

type clistruct struct {
	logging.LoggingConfig

	SaveDefaultConfig bool `help:"Save default config and exit"`

	CwdMode                string   `help:"How to display the current directory" default:"fancy"`
	CwdMaxDepth            int      `help:"Maximum number of directories to show in path" default:"5"`
	CwdMaxDirSize          int      `help:"Maximum number of letters displayed for each directory in the path" default:"-1"`
	ColorizeHostname       bool     `help:"Colorize the hostname based on a hash of itself, or use the PLGO_HOSTNAMEFG and PLGO_HOSTNAMEBG env vars (both need to be set)." default:"false"`
	HostnameOnlyIfSSH      bool     `help:"Show hostname only for SSH connections" default:"false"`
	SshAlternateIcon       bool     `help:"Show the older, original icon for SSH connections" default:"false"`
	EastAsianWidth         bool     `help:"Use East Asian Ambiguous Widths" default:"false"`
	Newline                bool     `help:"Show the prompt on a new line" default:"false"`
	StaticPromptIndicator  bool     `help:"Always show the prompt indicator with the default color, never with the error color" default:"false"`
	VenvNameSizeLimit      int      `help:"Show indicator instead of virtualenv name if name is longer than this limit (defaults to 0, which is unlimited)" default:"0"`
	Jobs                   int      `help:"Number of jobs currently running" default:"0"`
	GitAssumeUnchangedSize int64    `help:"Disable checking for changed/edited files in git repositories where the index is larger than this size (in KB), improves performance" default:"2048"`
	GitDisableStats        []string `help:"Comma-separated list to disable individual git statuses, (valid choices: ahead, behind, staged, notStaged, untracked, conflicted, stashed)"`
	GitMode                string   `help:"How to display git status, (valid choices: fancy, compact, simple)" default:"fancy"`
	Mode                   string   `help:"The characters used to make separators between segments, (valid choices: patched, compatible, flat)" default:"patched"`
	Theme                  string   `help:"Set this to the theme you want to use, (valid choices: default, low-contrast, gruvbox, solarized-dark16, solarized-light16)" default:"default"`
	Shell                  string   `help:"overwrite the shell (valid choices: autodetect, bare, bash, zsh)" enum:"autodetect,bare,bash,zsh" default:"autodetect"`
	Modules                []string `help:"The list of modules to load, separated by ','. Unrecognized modules will be invoked as 'powerline-go-MODULE' executable plugins and should output a (possibly empty) list of JSON objects that unmarshal to powerline-go's Segment structs." default:"venv,user,host,ssh,cwd,perms,git,hg,jobs,exit,root"`
	ModulesRight           []string `help:"The list of modules to load anchored to the right, separated by ','. Unrecognized modules will be invoked as 'powerline-go-MODULE' executable plugins and should output a (possibly empty) list of JSON objects that unmarshal to powerline-go's Segment structs."`
	Priority               []string `help:"Segments sorted by priority, if not enough space exists, the least priorized segments are removed first. Separate with ','" default:"root,cwd,user,host,ssh,perms,git-branch,git-status,hg,jobs,exit,cwd-path"`
	MaxWidthPercentage     int      `help:"Maximum width of the shell that the prompt may use, in percent. Setting this to 0 disables the shrinking subsystem." default:"0"`
	TruncateSegmentWidth   int      `help:"Maximum width of a segment, segments longer than this will be shortened if space is limited. Setting this to 0 disables it." default:"16"`
	PrevError              int      `help:"Exit code of previously executed command" default:"0"`
	NumericExitCodes       bool     `help:"Shows numeric exit codes for errors." default:"false"`
	IgnoreRepos            []string `help:"A list of git repos to ignore. Separate with ','. Repos are identified by their root directory."`
	ShortenGKENames        bool     `help:"Shortens names for GKE Kube clusters." default:"false"`
	ShortenEKSNames        bool     `help:"Shortens names for EKS Kube clusters." default:"false"`
	ShortenOpenshiftNames  bool     `help:"Shortens names for Openshift Kube clusters." default:"false"`
	ShellVar               string   `help:"A shell variable to add to the segments."`
	ShellVarNoWarnEmpty    bool     `help:"Disables warning for empty shell variable." default:"false"`
	TrimADDomain           bool     `helP:"Trim the Domainname from the AD username." default:"false"`
	Duration               string   `help:"The elapsed clock-time of the previous command"`
	DurationMin            string   `help:"The minimal time a command has to take before the duration segment is shown" default:"0"`
	DurationLowPrecision   bool     `help:"Use low precision timing for duration with milliseconds as maximum resolution" default:"false"`
	Eval                   bool     `help:"Output prompt in 'eval' format." default:"false"`
	Condensed              bool     `help:"Remove spacing between segments" default:"false"`
	Time                   string   `help:"The layout string how a reference time should be represented. The reference time is predefined and not user choosen. Consult the golang documentation for details: https://pkg.go.dev/time#example-Time.Format" default:"15:04:05"`
	ViMode                 string   `help:"The current vi-mode (eg. KEYMAP for zsh) for vi-module module"`

	PathAliases map[string]string `help:"One or more aliases from a path to a short name. Separate with ','. An alias maps a path like foo/bar/baz to a short name like FBB. Specify these as key/value pairs like foo/bar/baz=FBB. Use '~' for your home dir. You may need to escape this character to avoid shell substitution." mapsep:","`
}

var cli clistruct

func main() {
	ctx := kong.Parse(&cli, kong.UsageOnError(), kong.Configuration(kong.JSON, config.ConfigPath()))

	cli.Json = false
	logging.Setup(&cli.LoggingConfig)

	cfg := config.New()
	updateConfFromCLi(&cfg, cli)

	if cfg.SaveDefaultConfig {
		if err := cfg.EnsureExists(); err != nil {
			log.Fatal().Msg("could not write default config")
		}
		ctx.Exit(0)
	}

	if strings.HasSuffix(cfg.Theme, ".json") {
		file, err := os.ReadFile(cfg.Theme)
		if err == nil {
			theme := cfg.Themes["default"]
			err = json.Unmarshal(file, &theme)
			if err == nil {
				cfg.Themes[cfg.Theme] = theme
			} else {
				log.Error().Err(err).Msg("error reading theme")
			}
		}
	}

	if strings.HasSuffix(cfg.Mode, ".json") {
		file, err := os.ReadFile(cfg.Mode)
		if err == nil {
			symbols := cfg.Modes["compatible"]
			err = json.Unmarshal(file, &symbols)
			if err == nil {
				cfg.Modes[cfg.Mode] = symbols
			} else {
				log.Error().Err(err).Msg("error reading mode")
			}
		}
	}

	p := pwl.NewPowerline(cfg, config.AlignLeft)
	if p.SupportsRightModules() && p.HasRightModules() && !cfg.Eval {
		log.Fatal().Msg("'--modules-right' requires '--eval' mode")
	}

	fmt.Print(p.Render())

	ctx.Exit(0)
}

func updateConfFromCLi(cfg *config.Config, cli clistruct) {
	cfg.CwdMode = cli.CwdMode
	cfg.CwdMaxDepth = cli.CwdMaxDepth
	cfg.CwdMaxDirSize = cli.CwdMaxDirSize
	cfg.ColorizeHostname = cli.ColorizeHostname
	cfg.HostnameOnlyIfSSH = cli.HostnameOnlyIfSSH
	cfg.SshAlternateIcon = cli.SshAlternateIcon
	cfg.EastAsianWidth = cli.EastAsianWidth
	cfg.PromptOnNewLine = cli.Newline
	cfg.StaticPromptIndicator = cli.StaticPromptIndicator
	cfg.VenvNameSizeLimit = cli.VenvNameSizeLimit
	cfg.Jobs = cli.Jobs
	cfg.GitAssumeUnchangedSize = cli.GitAssumeUnchangedSize
	if len(cli.GitDisableStats) > 0 {
		cfg.GitDisableStats = cli.GitDisableStats
	}
	cfg.GitMode = cli.GitMode
	cfg.Mode = cli.Mode
	cfg.Theme = cli.Theme
	if cli.Shell != "autodetect" {
		cfg.Shell = cli.Shell
	}
	cfg.Modules = cli.Modules
	cfg.ModulesRight = cli.ModulesRight
	cfg.Priority = cli.Priority
	cfg.MaxWidthPercentage = cli.MaxWidthPercentage
	cfg.TruncateSegmentWidth = cli.TruncateSegmentWidth
	cfg.PrevError = cli.PrevError
	cfg.NumericExitCodes = cli.NumericExitCodes
	if len(cli.IgnoreRepos) > 0 {
		cfg.IgnoreRepos = cli.IgnoreRepos
	}
	cfg.ShortenGKENames = cli.ShortenGKENames
	cfg.ShortenEKSNames = cli.ShortenEKSNames
	cfg.ShortenOpenshiftNames = cli.ShortenOpenshiftNames
	cfg.ShellVar = cli.ShellVar
	cfg.ShellVarNoWarnEmpty = cli.ShellVarNoWarnEmpty
	cfg.TrimADDomain = cli.TrimADDomain
	cfg.Duration = cli.Duration
	cfg.DurationMin = cli.DurationMin
	cfg.DurationLowPrecision = cli.DurationLowPrecision
	cfg.Eval = cli.Eval
	cfg.Condensed = cli.Condensed
	cfg.Time = cli.Time
	cfg.ViMode = cli.ViMode
	cfg.PathAliases = cli.PathAliases
}
