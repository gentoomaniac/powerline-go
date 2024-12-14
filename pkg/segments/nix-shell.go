package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func NixShell(cfg config.Config, align config.Alignment) []Segment {
	var nixShell string
	nixShell, _ = os.LookupEnv("IN_NIX_SHELL")
	if nixShell == "" {
		return []Segment{}
	}
	return []Segment{{
		Name:       "nix-shell",
		Content:    "\uf313",
		Foreground: cfg.SelectedTheme().NixShellFg,
		Background: cfg.SelectedTheme().NixShellBg,
	}}
}
