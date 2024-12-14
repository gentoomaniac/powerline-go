//go:build broken

package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func SSH(theme config.Theme) []Segment {
	sshClient, _ := os.LookupEnv("SSH_CLIENT")
	if sshClient == "" {
		return []Segment{}
	}
	var networkIcon string
	if p.cfg.SshAlternateIcon {
		networkIcon = p.symbols.NetworkAlternate
	} else {
		networkIcon = p.symbols.Network
	}

	return []Segment{{
		Name:       "ssh",
		Content:    networkIcon,
		Foreground: theme.SSHFg,
		Background: theme.SSHBg,
	}}
}
