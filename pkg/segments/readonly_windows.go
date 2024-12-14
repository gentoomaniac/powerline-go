//go:build windows
// +build windows

package segments

import (
	"os"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

func Perms(cfg config.Config) []Segment {
	cwd := getValidCwd()
	const W_USR = 0o002
	// Check user's permissions on directory in a portable but probably slower way
	fileInfo, _ := os.Stat(cwd)
	if fileInfo.Mode()&W_USR == W_USR {
		return []Segment{}
	}
	return []Segment{{
		Name:       "perms",
		Content:    config.Defaults.Symbols().Lock,
		Foreground: config.Defaults.SelectedTheme().ReadonlyFg,
		Background: config.Defaults.SelectedTheme().ReadonlyBg,
	}}
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
