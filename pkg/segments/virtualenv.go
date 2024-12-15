package segments

import (
	"os"
	"path"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	"gopkg.in/ini.v1"

	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func VirtualEnv(theme config.Theme) []segment {
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
		return []segment{}
	}
	envName := path.Base(env)
	if p.cfg.VenvNameSizeLimit > 0 && len(envName) > p.cfg.VenvNameSizeLimit {
		envName = p.symbols.VenvIndicator
	}

	return []segment{{
		Name:       "venv",
		Content:    escapeVariables(p, envName),
		Foreground: theme.VirtualEnvFg,
		Background: theme.VirtualEnvBg,
	}}
}
