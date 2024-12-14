package segments

import (
	"io/ioutil"
	"os"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

const wsFile = "./.terraform/environment"

func TerraformWorkspace(theme config.Theme) []segment {
	stat, err := os.Stat(wsFile)
	if err != nil {
		return []segment{}
	}
	if stat.IsDir() {
		return []segment{}
	}
	workspace, err := ioutil.ReadFile(wsFile)
	if err != nil {
		return []segment{}
	}
	return []segment{{
		Name:       "terraform-workspace",
		Content:    string(workspace),
		Foreground: theme.TFWsFg,
		Background: theme.TFWsBg,
	}}
}
