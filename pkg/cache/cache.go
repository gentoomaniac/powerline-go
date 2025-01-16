package cache

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

var DefaultCachePath = []string{".cache", "powerline-go"}

func Save(entryName string, content []byte) error {
	home, _ := os.UserHomeDir()
	p := []string{home}
	path := filepath.Join(append(p, DefaultCachePath...)...)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(filepath.Join(path, entryName), content, 0o600)
}

func Read(entryName string) ([]byte, error) {
	home, _ := os.UserHomeDir()
	p := []string{home}
	path := filepath.Join(append(p, DefaultCachePath...)...)

	return os.ReadFile(filepath.Join(path, entryName))
}

func Age(entryName string) time.Time {
	home, _ := os.UserHomeDir()
	p := []string{home}
	path := filepath.Join(append(p, DefaultCachePath...)...)

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Error().Err(err).Msg("geting age of cached object failed")
		return time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	}

	return fileInfo.ModTime()
}
