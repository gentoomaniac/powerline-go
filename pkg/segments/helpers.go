package segments

import (
	"os"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

// TODO: this is a shitty place, needs better structure
func escapeVariables(cfg config.Config, pathSegment string) string {
	pathSegment = strings.Replace(pathSegment, `\`, cfg.CurrentShell().EscapedBackslash, -1)
	pathSegment = strings.Replace(pathSegment, "`", cfg.CurrentShell().EscapedBacktick, -1)
	pathSegment = strings.Replace(pathSegment, `$`, cfg.CurrentShell().EscapedDollar, -1)
	return pathSegment
}

func getValidCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		var exists bool
		cwd, exists = os.LookupEnv("PWD")
		if !exists {
			log.Warn().Msg("Your current directory is invalid.")
			print("> ")
			os.Exit(1)
		}
	}

	parts := strings.Split(cwd, string(os.PathSeparator))
	up := cwd

	for len(parts) > 0 && !pathExists(up) {
		parts = parts[:len(parts)-1]
		up = strings.Join(parts, string(os.PathSeparator))
	}
	if cwd != up {
		log.Warn().Msgf("Your current directory is invalid. Lowest valid directory: %s", up)
	}
	return cwd
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
