package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func SSH(theme config.Theme) []segment {
	sshClient, _ := os.LookupEnv("SSH_CLIENT")
	if sshClient == "" {
		return []segment{}
	}
	var networkIcon string
	if p.cfg.SshAlternateIcon {
		networkIcon = p.symbols.NetworkAlternate
	} else {
		networkIcon = p.symbols.Network
	}

	return []segment{{
		Name:       "ssh",
		Content:    networkIcon,
		Foreground: theme.SSHFg,
		Background: theme.SSHBg,
	}}
}
