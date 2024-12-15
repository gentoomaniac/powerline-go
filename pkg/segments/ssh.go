package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func SSH(cfg config.Config, align config.Alignment) []Segment {
	sshClient, _ := os.LookupEnv("SSH_CLIENT")
	if sshClient == "" {
		return []Segment{}
	}
	var networkIcon string
	if cfg.SshAlternateIcon {
		networkIcon = cfg.Symbols().NetworkAlternate
	} else {
		networkIcon = cfg.Symbols().Network
	}

	return []Segment{{
		Name:       "ssh",
		Content:    networkIcon,
		Foreground: cfg.Theme.SSHFg,
		Background: cfg.Theme.SSHBg,
	}}
}
