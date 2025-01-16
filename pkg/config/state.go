package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/gentoomaniac/powerline-go/pkg/shellinfo"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/process"
)

var (
	DefaultConfigPath = []string{".config", "powerline-go"}
	ConfigFileName    = "config.json"
)

type State struct {
	CwdMode                string
	CwdMaxDepth            int
	CwdMaxDirSize          int
	ColorizeHostname       bool
	HostnameOnlyIfSsh      bool
	SshAlternateIcon       bool
	EastAsianWidth         bool
	PromptOnNewLine        bool
	StaticPromptIndicator  bool
	VenvNameSizeLimit      int
	Jobs                   int
	GitAssumeUnchangedSize int64
	GitDisableStats        []string
	GitMode                string
	GithubToken            string
	Mode                   string
	Theme                  Theme
	Shell                  string
	Modules                []string
	ModulesRight           []string
	Priority               []string
	MaxWidthPercentage     int
	TruncateSegmentWidth   int
	PrevError              int
	NumericExitCodes       bool
	IgnoreRepos            []string
	ShortenGkeNames        bool
	ShortenEksNames        bool
	ShortenOpenshiftNames  bool
	ShellVar               string
	ShellVarNoWarnEmpty    bool
	TrimAdDomain           bool
	Duration               string
	DurationMin            string
	DurationLowPrecision   bool
	Eval                   bool
	Condensed              bool
	TimeFormat             string
	ViMode                 string
	PathAliases            map[string]string

	Modes     map[string]SymbolTemplate      `json:"-"`
	Shells    map[string]shellinfo.ShellInfo `json:"-"`
	Userinfo  *user.User                     `json:"-"`
	Hostname  string                         `json:"-"`
	Cwd       string                         `json:"-"`
	ThemeName string                         `json:"-"`
}

func NewStateFromConfig(cfg Config) State {
	var shellExe string
	proc, err := process.NewProcess(int32(os.Getppid()))
	if err == nil {
		shellExe, _ = proc.Exe()
	}
	if shellExe == "" {
		shellExe = os.Getenv("SHELL")
	}
	shell := detectShell(shellExe)

	user, err := user.Current()
	if err != nil {
		log.Error().Err(err).Msg("failed getting userinfo")
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Error().Err(err).Msg("couldn't determine hostname")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("could not determine current working directory")
	}

	defaults.Userinfo = user
	defaults.Hostname = hostname
	defaults.Shell = shell
	defaults.Cwd = cwd

	defaults.CwdMode = cfg.CwdMode
	defaults.CwdMaxDepth = cfg.CwdMaxDepth
	defaults.CwdMaxDirSize = cfg.CwdMaxDirSize
	defaults.ColorizeHostname = cfg.ColorizeHostname
	defaults.HostnameOnlyIfSsh = cfg.HostnameOnlyIfSsh
	defaults.SshAlternateIcon = cfg.SshAlternateIcon
	defaults.EastAsianWidth = cfg.EastAsianWidth
	defaults.PromptOnNewLine = cfg.Newline
	defaults.StaticPromptIndicator = cfg.StaticPromptIndicator
	defaults.VenvNameSizeLimit = cfg.VenvNameSizeLimit
	defaults.Jobs = cfg.Jobs
	defaults.GitAssumeUnchangedSize = cfg.GitAssumeUnchangedSize
	if len(cfg.GitDisableStats) > 0 {
		defaults.GitDisableStats = cfg.GitDisableStats
	}
	defaults.GitMode = cfg.GitMode
	defaults.GithubToken = cfg.GithubToken
	defaults.Mode = cfg.Mode
	if cfg.Shell != "autodetect" {
		defaults.Shell = cfg.Shell
	}
	defaults.Modules = cfg.Modules
	defaults.ModulesRight = cfg.ModulesRight
	defaults.Priority = cfg.Priority
	defaults.MaxWidthPercentage = cfg.MaxWidthPercentage
	defaults.TruncateSegmentWidth = cfg.TruncateSegmentWidth
	defaults.PrevError = cfg.PrevError
	defaults.NumericExitCodes = cfg.NumericExitCodes
	if len(cfg.IgnoreRepos) > 0 {
		defaults.IgnoreRepos = cfg.IgnoreRepos
	}
	defaults.ShortenGkeNames = cfg.ShortenGkeNames
	defaults.ShortenEksNames = cfg.ShortenEksNames
	defaults.ShortenOpenshiftNames = cfg.ShortenOpenshiftNames
	defaults.ShellVar = cfg.ShellVar
	defaults.ShellVarNoWarnEmpty = cfg.ShellVarNoWarnEmpty
	defaults.TrimAdDomain = cfg.TrimAdDomain
	defaults.Duration = cfg.Duration
	defaults.DurationMin = cfg.DurationMin
	defaults.DurationLowPrecision = cfg.DurationLowPrecision
	defaults.Eval = cfg.Eval
	defaults.Condensed = cfg.Condensed
	defaults.TimeFormat = cfg.TimeFormat
	defaults.ViMode = cfg.ViMode
	defaults.PathAliases = cfg.PathAliases
	defaults.ThemeName = cfg.Theme
	return defaults
}

func (cfg *State) Symbols() SymbolTemplate {
	return cfg.Modes[cfg.Mode]
}

func (cfg *State) CurrentShell() shellinfo.ShellInfo {
	return cfg.Shells[cfg.Shell]
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	path := []string{home}
	path = append(path, DefaultConfigPath...)
	return filepath.Join(path...)
}
