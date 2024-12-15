package segments

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func Direnv(cfg config.Config, align config.Alignment) []Segment {
	content := os.Getenv("DIRENV_DIR")
	if content == "" {
		return []Segment{}
	}
	if strings.TrimPrefix(content, "-") == cfg.Userinfo.HomeDir {
		content = "~"
	} else {
		content = filepath.Base(content)
	}

	return []Segment{{
		Name:       "direnv",
		Content:    content,
		Foreground: cfg.Theme.DotEnvFg,
		Background: cfg.Theme.DotEnvBg,
	}}
}
