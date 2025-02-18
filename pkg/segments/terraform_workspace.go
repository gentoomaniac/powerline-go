package segments

import (
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

const wsFile = "./.terraform/environment"

func TerraformWorkspace(cfg config.State, align config.Alignment) []Segment {
	stat, err := os.Stat(wsFile)
	if err != nil {
		return []Segment{}
	}
	if stat.IsDir() {
		return []Segment{}
	}
	workspace, err := os.ReadFile(wsFile)
	if err != nil {
		return []Segment{}
	}
	return []Segment{{
		Name:       "terraform-workspace",
		Content:    string(workspace),
		Foreground: cfg.Theme.TFWsFg,
		Background: cfg.Theme.TFWsBg,
	}}
}
