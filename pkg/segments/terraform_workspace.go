package segments

import (
	"io/ioutil"
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

const wsFile = "./.terraform/environment"

func TerraformWorkspace(theme config.Theme) []Segment {
	stat, err := os.Stat(wsFile)
	if err != nil {
		return []Segment{}
	}
	if stat.IsDir() {
		return []Segment{}
	}
	workspace, err := ioutil.ReadFile(wsFile)
	if err != nil {
		return []Segment{}
	}
	return []Segment{{
		Name:       "terraform-workspace",
		Content:    string(workspace),
		Foreground: theme.TFWsFg,
		Background: theme.TFWsBg,
	}}
}
