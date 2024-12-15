package segments

import (
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

// TODO: this is a shitty place, needs better structure
func escapeVariables(cfg config.Config, pathSegment string) string {
	pathSegment = strings.Replace(pathSegment, `\`, cfg.CurrentShell().EscapedBackslash, -1)
	pathSegment = strings.Replace(pathSegment, "`", cfg.CurrentShell().EscapedBacktick, -1)
	pathSegment = strings.Replace(pathSegment, `$`, cfg.CurrentShell().EscapedDollar, -1)
	return pathSegment
}
