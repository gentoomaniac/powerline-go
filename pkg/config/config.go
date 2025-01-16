package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gentoomaniac/powerline-go/pkg/cli"
	"github.com/gentoomaniac/powerline-go/pkg/logging"
)

type Config struct {
	logging.LoggingConfig `json:"-"`

	SaveConfig bool `help:"Save default config and exit" default:"false" json:"-"`

	CwdMode                string            `help:"How to display the current directory" default:"fancy" json:"cwdMode"`
	CwdMaxDepth            int               `help:"Maximum number of directories to show in path" default:"5" json:"cwdMaxDepth"`
	CwdMaxDirSize          int               `help:"Maximum number of letters displayed for each directory in the path" default:"-1" json:"cwdMaxDirSize"`
	ColorizeHostname       bool              `help:"Colorize the hostname based on a hash of itself, or use the PLGO_HOSTNAMEFG and PLGO_HOSTNAMEBG env vars (both need to be set)." default:"false" json:"colorizeHostname"`
	HostnameOnlyIfSsh      bool              `help:"Show hostname only for SSH connections" default:"false" json:"hostnameOnlyIfSsh"`
	SshAlternateIcon       bool              `help:"Show the older, original icon for SSH connections" default:"false" json:"sshAlternateIcon"`
	EastAsianWidth         bool              `help:"Use East Asian Ambiguous Widths" default:"false" json:"eastAsianWidth"`
	Newline                bool              `help:"Show the prompt on a new line" default:"false" json:"newline"`
	StaticPromptIndicator  bool              `help:"Always show the prompt indicator with the default color, never with the error color" default:"false" json:"staticPromptIndicator"`
	VenvNameSizeLimit      int               `help:"Show indicator instead of virtualenv name if name is longer than this limit (defaults to 0, which is unlimited)" default:"0" json:"venvNameSizeLimit"`
	Jobs                   int               `help:"Number of jobs currently running" default:"0" json:"jobs"`
	GitAssumeUnchangedSize int64             `help:"Disable checking for changed/edited files in git repositories where the index is larger than this size (in KB), improves performance" default:"2048" json:"gitAssumeUnchangedSize"`
	GitDisableStats        []string          `help:"Comma-separated list to disable individual git statuses, (valid choices: ahead, behind, staged, notStaged, untracked, conflicted, stashed)" json:"gitDisableStats"`
	GitMode                string            `help:"How to display git status, (valid choices: fancy, compact, simple)" default:"fancy" json:"gitMode"`
	GithubToken            string            `help:"Github PAT (classic) that has access to notifications" default:"fancy" env:"GH_TOKEN" json:"ghToken"`
	Mode                   string            `help:"The characters used to make separators between segments, (valid choices: patched, compatible, flat)" default:"patched" json:"mode"`
	Theme                  string            `help:"Set this to the theme you want to use, (valid choices: default, low-contrast, gruvbox, solarized-dark16, solarized-light16)" default:"default" json:"theme"`
	Shell                  string            `help:"overwrite the shell (valid choices: autodetect, bare, bash, zsh)" enum:"autodetect,bare,bash,zsh" default:"autodetect" json:"shell"`
	Modules                []string          `help:"The list of modules to load, separated by ','. Unrecognized modules will be invoked as 'powerline-go-MODULE' executable plugins and should output a (possibly empty) list of JSON objects that unmarshal to powerline-go's Segment structs." default:"venv,user,host,ssh,cwd,perms,git,hg,jobs,exit,root" json:"modules"`
	ModulesRight           []string          `help:"The list of modules to load anchored to the right, separated by ','. Unrecognized modules will be invoked as 'powerline-go-MODULE' executable plugins and should output a (possibly empty) list of JSON objects that unmarshal to powerline-go's Segment structs." json:"modulesRight"`
	Priority               []string          `help:"Segments sorted by priority, if not enough space exists, the least priorized segments are removed first. Separate with ','" default:"root,cwd,user,host,ssh,perms,git-branch,git-status,hg,jobs,exit,cwd-path" json:"priority"`
	MaxWidthPercentage     int               `help:"Maximum width of the shell that the prompt may use, in percent. Setting this to 0 disables the shrinking subsystem." default:"0" json:"maxWidthPercentage"`
	TruncateSegmentWidth   int               `help:"Maximum width of a segment, segments longer than this will be shortened if space is limited. Setting this to 0 disables it." default:"16" json:"truncateSegmentWidth"`
	PrevError              int               `help:"Exit code of previously executed command" default:"0" json:"prevError"`
	NumericExitCodes       bool              `help:"Shows numeric exit codes for errors." default:"false" json:"numericExitCodes"`
	IgnoreRepos            []string          `help:"A list of git repos to ignore. Separate with ','. Repos are identified by their root directory." json:"ignoreRepos"`
	ShortenGkeNames        bool              `help:"Shortens names for GKE Kube clusters." default:"false" json:"shortenGKENames"`
	ShortenEksNames        bool              `help:"Shortens names for EKS Kube clusters." default:"false" json:"shortenEKSNames"`
	ShortenOpenshiftNames  bool              `help:"Shortens names for Openshift Kube clusters." default:"false" json:"shortenOpenshiftNames"`
	ShellVar               string            `help:"A shell variable to add to the segments." json:"shellVar"`
	ShellVarNoWarnEmpty    bool              `help:"Disables warning for empty shell variable." default:"false" json:"shellVarNoWarnEmpty"`
	TrimAdDomain           bool              `helP:"Trim the Domainname from the AD username." default:"false" json:"trimAdDomain"`
	Duration               string            `help:"The elapsed clock-time of the previous command" json:"duration"`
	DurationMin            string            `help:"The minimal time a command has to take before the duration segment is shown" default:"0" json:"durationMin"`
	DurationLowPrecision   bool              `help:"Use low precision timing for duration with milliseconds as maximum resolution" default:"false" json:"durationLowPrecision"`
	Eval                   bool              `help:"Output prompt in 'eval' format." default:"false" json:"eval"`
	Condensed              bool              `help:"Remove spacing between segments" default:"false" json:"condensed"`
	TimeFormat             string            `help:"The layout string how a reference time should be represented. The reference time is predefined and not user choosen. Consult the golang documentation for details: https://pkg.go.dev/time#example-Time.Format" default:"15:04:05" json:"timeFormat"`
	ViMode                 string            `help:"The current vi-mode (eg. KEYMAP for zsh) for vi-module module" json:"viMode"`
	PathAliases            map[string]string `help:"One or more aliases from a path to a short name. Separate with ','. An alias maps a path like foo/bar/baz to a short name like FBB. Specify these as key/value pairs like foo/bar/baz=FBB. Use '~' for your home dir. You may need to escape this character to avoid shell substitution." mapsep:"," json:"pathAliases"`

	Version cli.VersionFlag `short:"V" help:"Display version."`
}

func (conf Config) save() error {
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(ConfigPath(), ConfigFileName), data, 0o644)
}

func (conf Config) Save() error {
	home, _ := os.UserHomeDir()
	p := []string{home}
	path := filepath.Join(append(p, DefaultConfigPath...)...)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return conf.save()
}
