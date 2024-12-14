package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func NixShell(theme config.Theme) []Segment {
	var nixShell string
	nixShell, _ = os.LookupEnv("IN_NIX_SHELL")
	if nixShell == "" {
		return []Segment{}
	}
	return []Segment{{
		Name:       "nix-shell",
		Content:    "\uf313",
		Foreground: theme.NixShellFg,
		Background: theme.NixShellBg,
	}}
}
