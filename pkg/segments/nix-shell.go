package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func NixShell(cfg config.State, align config.Alignment) []Segment {
	var nixShell string
	nixShell, _ = os.LookupEnv("IN_NIX_SHELL")
	if nixShell == "" {
		return []Segment{}
	}
	return []Segment{{
		Name:       "nix-shell",
		Content:    "\uf313",
		Foreground: cfg.Theme.NixShellFg,
		Background: cfg.Theme.NixShellBg,
	}}
}
