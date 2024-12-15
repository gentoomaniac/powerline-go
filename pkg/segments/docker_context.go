package segments

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

type DockerContextConfig struct {
	CurrentContext string `json:"currentContext"`
}

func DockerContext(cfg config.Config, align config.Alignment) []Segment {
	context := "default"
	home, _ := os.LookupEnv("HOME")
	contextFolder := filepath.Join(home, ".docker", "contexts")
	configFile := filepath.Join(home, ".docker", "config.json")
	contextEnvVar := os.Getenv("DOCKER_CONTEXT")

	if contextEnvVar != "" {
		context = contextEnvVar
	} else {
		stat, err := os.Stat(contextFolder)
		if err == nil && stat.IsDir() {
			dockerConfigFile, err := os.ReadFile(configFile)
			if err == nil {
				var dockerConfig DockerContextConfig
				err = json.Unmarshal(dockerConfigFile, &dockerConfig)
				if err == nil && dockerConfig.CurrentContext != "" {
					context = dockerConfig.CurrentContext
				}
			}
		}
	}

	// Don’t show the default context
	if context == "default" {
		return []Segment{}
	}

	return []Segment{{
		Name:       "docker-context",
		Content:    "🐳" + context,
		Foreground: cfg.Theme.PlEnvFg,
		Background: cfg.Theme.PlEnvBg,
	}}
}
