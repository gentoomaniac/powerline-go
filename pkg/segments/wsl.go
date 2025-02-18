package segments

import (
	"net/url"
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func WSL(cfg config.State, align config.Alignment) []Segment {
	var WSL string
	WSLMachineName, _ := os.LookupEnv("WSL_DISTRO_NAME")
	WSLHost, _ := os.LookupEnv("NAME")

	if WSLMachineName != "" {
		WSL = WSLMachineName
	} else if WSLHost != " " {
		u, err := url.Parse(WSLHost)
		if err == nil {
			WSL = u.Host
		}
	}

	if WSL != "" {
		return []Segment{{
			Name:       "WSL",
			Content:    WSL,
			Foreground: cfg.Theme.WSLMachineFg,
			Background: cfg.Theme.WSLMachineBg,
		}}
	}
	return []Segment{}
}
