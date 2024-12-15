package segments

import (
	"os"
	"path"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	"gopkg.in/ini.v1"
)

func VirtualEnv(cfg config.Config, align config.Alignment) []Segment {
	env := os.Getenv("VIRTUAL_ENV_PROMPT")
	if strings.HasPrefix(env, "(") && strings.HasSuffix(env, ") ") {
		env = strings.TrimPrefix(env, "(")
		env = strings.TrimSuffix(env, ") ")
	}
	if env == "" {
		venv, _ := os.LookupEnv("VIRTUAL_ENV")
		if venv != "" {
			cfg, err := ini.Load(path.Join(venv, "pyvenv.cfg"))
			if err == nil {
				// python >= 3.6 the venv module will not insert a prompt
				// key unless the `--prompt` flag is passed to the module
				// or if calling with the prompt arg EnvBuilder
				// otherwise env evaluates to an empty string, per return
				// of ini.File.Section.Key
				if pyEnv := cfg.Section("").Key("prompt").String(); pyEnv != "" {
					env = pyEnv
				}
			}
			if env == "" {
				env = venv
			}
		}
	}
	if env == "" {
		env, _ = os.LookupEnv("CONDA_ENV_PATH")
	}
	if env == "" {
		env, _ = os.LookupEnv("CONDA_DEFAULT_ENV")
	}
	if env == "" {
		env, _ = os.LookupEnv("PYENV_VERSION")
	}
	if env == "" {
		return []Segment{}
	}
	envName := path.Base(env)
	if cfg.VenvNameSizeLimit > 0 && len(envName) > cfg.VenvNameSizeLimit {
		envName = cfg.Symbols().VenvIndicator
	}

	return []Segment{{
		Name:       "venv",
		Content:    escapeVariables(cfg, envName),
		Foreground: cfg.Theme.VirtualEnvFg,
		Background: cfg.Theme.VirtualEnvBg,
	}}
}
