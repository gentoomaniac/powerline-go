package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func NixShell(theme config.Theme) []segment {
	var nixShell string
	nixShell, _ = os.LookupEnv("IN_NIX_SHELL")
	if nixShell == "" {
		return []segment{}
	}
	return []segment{{
		Name:       "nix-shell",
		Content:    "\uf313",
		Foreground: theme.NixShellFg,
		Background: theme.NixShellBg,
	}}
}
