package segments

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func Direnv(theme config.Theme) []segment {
	content := os.Getenv("DIRENV_DIR")
	if content == "" {
		return []segment{}
	}
	if strings.TrimPrefix(content, "-") == p.userInfo.HomeDir {
		content = "~"
	} else {
		content = filepath.Base(content)
	}

	return []segment{{
		Name:       "direnv",
		Content:    content,
		Foreground: theme.DotEnvFg,
		Background: theme.DotEnvBg,
	}}
}
