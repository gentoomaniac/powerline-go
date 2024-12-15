package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gentoomaniac/powerline-go/pkg/shellinfo"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	defaultConfigPath = []string{".config", "powerline-go"}
	ConfigFileName    = "config.json"
)

type Config struct {
	CwdMode                string
	CwdMaxDepth            int
	CwdMaxDirSize          int
	ColorizeHostname       bool
	HostnameOnlyIfSSH      bool
	SshAlternateIcon       bool
	EastAsianWidth         bool
	PromptOnNewLine        bool
	StaticPromptIndicator  bool
	VenvNameSizeLimit      int
	Jobs                   int
	GitAssumeUnchangedSize int64
	GitDisableStats        []string
	GitMode                string
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
	ShortenGKENames        bool
	ShortenEKSNames        bool
	ShortenOpenshiftNames  bool
	ShellVar               string
	ShellVarNoWarnEmpty    bool
	TrimADDomain           bool
	Duration               string
	DurationMin            string
	DurationLowPrecision   bool
	Eval                   bool
	Condensed              bool
	Time                   string
	ViMode                 string
	PathAliases            map[string]string

	Modes     map[string]SymbolTemplate      `json:"-"`
	Shells    map[string]shellinfo.ShellInfo `json:"-"`
	Userinfo  *user.User                     `json:"-"`
	Hostname  string                         `json:"-"`
	Cwd       string                         `json:"-"`
	ThemeName string                         `json:"-"`
}

func New() Config {
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
	return defaults
}

func (cfg *Config) Symbols() SymbolTemplate {
	return cfg.Modes[cfg.Mode]
}

func (cfg *Config) CurrentShell() shellinfo.ShellInfo {
	return cfg.Shells[cfg.Shell]
}

func (cfg *Config) save() error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(ConfigPath(), ConfigFileName), data, 0o644)
}

func (cfg *Config) EnsureExists() error {
	home, _ := os.UserHomeDir()
	p := []string{home}
	path := filepath.Join(append(p, defaultConfigPath...)...)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("config already exists")
	}

	return cfg.save()
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	path := []string{home}
	path = append(path, defaultConfigPath...)
	return filepath.Join(path...)
}
