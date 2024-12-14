//go:build broken

package segments

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Direnv(theme config.Theme) []Segment {
	content := os.Getenv("DIRENV_DIR")
	if content == "" {
		return []Segment{}
	}
	if strings.TrimPrefix(content, "-") == p.userInfo.HomeDir {
		content = "~"
	} else {
		content = filepath.Base(content)
	}

	return []Segment{{
		Name:       "direnv",
		Content:    content,
		Foreground: theme.DotEnvFg,
		Background: theme.DotEnvBg,
	}}
}
