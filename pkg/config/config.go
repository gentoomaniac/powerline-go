package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gentoomaniac/powerline-go/pkg/logging"
	"github.com/gentoomaniac/powerline-go/pkg/shellinfo"
)

var defaultConfigPath = []string{".config", "powerline-go", "config.json"}

type (
	ShellMap map[string]shellinfo.ShellInfo
	ThemeMap map[string]Theme
)

type Config struct {
	logging.LoggingConfig

	SaveDefaultConfig bool `help:"Save default config and exit"`

	CwdMode                string   `help:"How to display the current directory"`
	CwdMaxDepth            int      `help:"Maximum number of directories to show in path"`
	CwdMaxDirSize          int      `help:"Maximum number of letters displayed for each directory in the path"`
	ColorizeHostname       bool     `help:"Colorize the hostname based on a hash of itself, or use the PLGO_HOSTNAMEFG and PLGO_HOSTNAMEBG env vars (both need to be set)."`
	HostnameOnlyIfSSH      bool     `help:"Show hostname only for SSH connections"`
	SshAlternateIcon       bool     `help:"Show the older, original icon for SSH connections"`
	EastAsianWidth         bool     `help:"Use East Asian Ambiguous Widths"`
	PromptOnNewLine        bool     `help:"Show the prompt on a new line"`
	StaticPromptIndicator  bool     `help:"Always show the prompt indicator with the default color, never with the error color"`
	VenvNameSizeLimit      int      `help:"Show indicator instead of virtualenv name if name is longer than this limit (defaults to 0, which is unlimited)"`
	Jobs                   int      `help:"Number of jobs currently running"`
	GitAssumeUnchangedSize int64    `help:"Disable checking for changed/edited files in git repositories where the index is larger than this size (in KB), improves performance"`
	GitDisableStats        []string `help:"Comma-separated list to disable individual git statuses, (valid choices: ahead, behind, staged, notStaged, untracked, conflicted, stashed)"`
	GitMode                string   `help:"How to display git status, (valid choices: fancy, compact, simple)"`
	Mode                   string   `help:"The characters used to make separators between segments, (valid choices: patched, compatible, flat)"`
	Theme                  string   `help:"Set this to the theme you want to use, (valid choices: default, low-contrast, gruvbox, solarized-dark16, solarized-light16)"`
	Shell                  string   `help:"Set this to your shell type, (valid choices: autodetect, bare, bash, zsh)"`
	Modules                []string `help:"The list of modules to load, separated by ','. Unrecognized modules will be invoked as 'powerline-go-MODULE' executable plugins and should output a (possibly empty) list of JSON objects that unmarshal to powerline-go's Segment structs." enum:"aws,bzr,cwd,direnv,docker,docker-context,dotenv,duration,exit,fossil,gcp,git,gitlite,goenv,hg,host,jobs,kube,load,newline,nix-shell,node,perlbrew,perms,plenv,rbenv,root,rvm,shell-var,shenv,ssh,svn,termtitle,terraform-workspace,time,user,venv,vgo,vi-mode,wsl"`
	ModulesRight           []string `help:"The list of modules to load anchored to the right, separated by ','. Unrecognized modules will be invoked as 'powerline-go-MODULE' executable plugins and should output a (possibly empty) list of JSON objects that unmarshal to powerline-go's Segment structs." enum:"aws,bzr,cwd,direnv,docker,docker-context,dotenv,duration,exit,fossil,gcp,git,gitlite,goenv,hg,host,jobs,kube,load,newline,nix-shell,node,perlbrew,perms,plenv,rbenv,root,rvm,shell-var,shenv,ssh,svn,termtitle,terraform-workspace,time,user,venv,vgo,vi-mode,wsl"`
	Priority               []string `help:"Segments sorted by priority, if not enough space exists, the least priorized segments are removed first. Separate with ','" enum:"aws,bzr,cwd,direnv,docker,docker-context,dotenv,duration,exit,fossil,gcp,git,gitlite,goenv,hg,host,jobs,kube,load,newline,nix-shell,node,perlbrew,perms,plenv,rbenv,root,rvm,shell-var,shenv,ssh,svn,termtitle,terraform-workspace,time,user,venv,vgo,vi-mode,wsl"`
	MaxWidthPercentage     int      `help:"Maximum width of the shell that the prompt may use, in percent. Setting this to 0 disables the shrinking subsystem."`
	TruncateSegmentWidth   int      `help:"Maximum width of a segment, segments longer than this will be shortened if space is limited. Setting this to 0 disables it."`
	PrevError              int      `help:"Exit code of previously executed command"`
	NumericExitCodes       bool     `help:"Shows numeric exit codes for errors."`
	IgnoreRepos            []string `help:"A list of git repos to ignore. Separate with ','. Repos are identified by their root directory."`
	ShortenGKENames        bool     `help:"Shortens names for GKE Kube clusters."`
	ShortenEKSNames        bool     `help:"Shortens names for EKS Kube clusters."`
	ShortenOpenshiftNames  bool     `help:"Shortens names for Openshift Kube clusters."`
	ShellVar               string   `help:"A shell variable to add to the segments."`
	ShellVarNoWarnEmpty    bool     `help:"Disables warning for empty shell variable."`
	TrimADDomain           bool     `helP:"Trim the Domainname from the AD username."`
	Duration               string   `help:"The elapsed clock-time of the previous command"`
	DurationMin            string   `help:"The minimal time a command has to take before the duration segment is shown"`
	DurationLowPrecision   bool     `help:"Use low precision timing for duration with milliseconds as maximum resolution"`
	Eval                   bool     `help:"Output prompt in 'eval' format."`
	Condensed              bool     `help:"Remove spacing between segments"`
	Time                   string   `help:"The layout string how a reference time should be represented. The reference time is predefined and not user choosen. Consult the golang documentation for details: https://pkg.go.dev/time#example-Time.Format"`
	ViMode                 string   `help:"The current vi-mode (eg. KEYMAP for zsh) for vi-module module"`

	PathAliases map[string]string `help:"One or more aliases from a path to a short name. Separate with ','. An alias maps a path like foo/bar/baz to a short name like FBB. Specify these as key/value pairs like foo/bar/baz=FBB. Use '~' for your home dir. You may need to escape this character to avoid shell substitution." mapsep:","`

	Modes  map[string]SymbolTemplate `hidden:""`
	Shells ShellMap                  `hidden:""`
	Themes ThemeMap                  `hidden:""`
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
