package segments

import (
	"net/url"
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Docker(theme config.Theme) []segment {
	var docker string
	dockerMachineName, _ := os.LookupEnv("DOCKER_MACHINE_NAME")
	dockerHost, _ := os.LookupEnv("DOCKER_HOST")

	if dockerMachineName != "" {
		docker = dockerMachineName
	} else if dockerHost != " " {
		u, err := url.Parse(dockerHost)
		if err == nil {
			docker = u.Host
		}
	}

	if docker == "" {
		return []segment{}
	}
	return []segment{{
		Name:       "docker",
		Content:    docker,
		Foreground: theme.DockerMachineFg,
		Background: theme.DockerMachineBg,
	}}
}
