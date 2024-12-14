package segments

import (
	"net/url"
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func WSL(theme config.Theme) []segment {
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
		return []segment{{
			Name:       "WSL",
			Content:    WSL,
			Foreground: theme.WSLMachineFg,
			Background: theme.WSLMachineBg,
		}}
	}
	return []segment{}
}
