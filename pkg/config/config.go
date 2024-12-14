package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gentoomaniac/powerline-go/pkg/logging"
	"github.com/gentoomaniac/powerline-go/pkg/shellinfo"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/process"
)

var defaultConfigPath = []string{".config", "powerline-go", "config.json"}

type (
	ShellMap map[string]shellinfo.ShellInfo
	ThemeMap map[string]Theme
)

type Config struct {
	logging.LoggingConfig

	SaveDefaultConfig bool

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
	Theme                  string
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

	PathAliases map[string]string

	Modes    map[string]SymbolTemplate
	Shells   ShellMap
	Themes   ThemeMap
	Userinfo *user.User
	Hostname string
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

	defaults.Userinfo = user
	defaults.Hostname = hostname
	defaults.Shell = shell
	return defaults
}

func (cfg *Config) SelectedTheme() Theme {
	return cfg.Themes[cfg.Theme]
}

func (cfg *Config) Symbols() SymbolTemplate {
	return cfg.Modes[cfg.Mode]
}

func (cfg *Config) CurrentShell() shellinfo.ShellInfo {
	return cfg.Shells[cfg.Shell]
}

func (cfg *Config) save() error {
	path := ConfigPath()
	tmp := cfg
	tmp.Themes = map[string]Theme{}
	tmp.Modes = map[string]SymbolTemplate{}
	tmp.Shells = map[string]shellinfo.ShellInfo{}
	data, err := json.MarshalIndent(tmp, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (cfg *Config) EnsureExists() error {
	home, _ := os.UserHomeDir()
	p := []string{home}
	path := filepath.Join(append(p, defaultConfigPath[:len(defaultConfigPath)-1]...)...)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return cfg.save()
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	path := []string{home}
	path = append(path, defaultConfigPath...)
	return filepath.Join(path...)
}
