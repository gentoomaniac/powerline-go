package segments

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

const (
	goenvVersionFileSuffix       = "/.go-version"
	goenvVersionEnvVar           = "GOENV_VERSION"
	goenvGlobalVersionFileSuffix = "/.goenv/version"
)

func runGoenvCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	out, err := command.Output()
	return string(out), err
}

// check GOENV_VERSION variable
func checkEnvForGoenvVersion() (string, error) {
	goenvVersion := os.Getenv(goenvVersionEnvVar)
	if len(goenvVersion) > 0 {
		return goenvVersion, nil
	} else {
		return "", fmt.Errorf("not found in %s", goenvVersionEnvVar)
	}
}

// check existence of .go-version in tree until root path
func checkForGoVersionFileInTree() (string, error) {
	var (
		workingDirectory string
		err              error
	)

	workingDirectory, err = os.Getwd()
	if err == nil {
		for workingDirectory != "/" {
			goVersion, goVersionErr := ioutil.ReadFile(workingDirectory + goenvVersionFileSuffix)
			if goVersionErr == nil {
				return strings.TrimSpace(string(goVersion)), nil
			}

			workingDirectory = filepath.Dir(workingDirectory)
		}
	}

	return "", fmt.Errorf("no %s file found in tree", goenvVersionFileSuffix)
}

// check for global version
func checkForGoenvGlobalVersion() (string, error) {
	homeDir, _ := os.UserHomeDir()
	globalGoVersion, err := os.ReadFile(homeDir + goenvGlobalVersionFileSuffix)
	if err == nil {
		return strings.TrimSpace(string(globalGoVersion)), nil
	} else {
		return "", fmt.Errorf("no global go version file found in %s", homeDir+goenvGlobalVersionFileSuffix)
	}
}

// retrieve goenv version output
func checkForGoenvOutput() (string, error) {
	// spawn goenv and print out version
	out, err := runGoenvCommand("goenv", "version")
	if err == nil {
		items := strings.Split(out, " ")
		if len(items) > 1 {
			return items[0], nil
		}
	}

	return "", fmt.Errorf("not found in goenv object")
}

func Goenv(cfg config.Config, align config.Alignment) []Segment {
	global, _ := checkForGoenvGlobalVersion()

	segment, err := checkEnvForGoenvVersion()
	if err != nil || segment == global {
		segment, err = checkForGoVersionFileInTree()
	}
	if err != nil || segment == global {
		segment, err = checkForGoenvOutput()
	}
	if err != nil || segment == global {
		return []Segment{}
	} else {
		return []Segment{{
			Name:       "goenv",
			Content:    segment,
			Foreground: cfg.SelectedTheme().GoenvFg,
			Background: cfg.SelectedTheme().GoenvBg,
		}}
	}
}
