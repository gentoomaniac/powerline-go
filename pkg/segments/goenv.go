package segments

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
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
		return "", fmt.Errorf("Not found in %s", goenvVersionEnvVar)
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

	return "", fmt.Errorf("No %s file found in tree", goenvVersionFileSuffix)
}

// check for global version
func checkForGoenvGlobalVersion() (string, error) {
	homeDir, _ := os.UserHomeDir()
	globalGoVersion, err := ioutil.ReadFile(homeDir + goenvGlobalVersionFileSuffix)
	if err == nil {
		return strings.TrimSpace(string(globalGoVersion)), nil
	} else {
		return "", fmt.Errorf("No global go version file found in %s", homeDir+goenvGlobalVersionFileSuffix)
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

	return "", fmt.Errorf("Not found in goenv object")
}

func Goenv(theme config.Theme) []segment {
	global, _ := checkForGoenvGlobalVersion()

	segment, err := checkEnvForGoenvVersion()
	if err != nil || segment == global {
		segment, err = checkForGoVersionFileInTree()
	}
	if err != nil || segment == global {
		segment, err = checkForGoenvOutput()
	}
	if err != nil || segment == global {
		return []segment{}
	} else {
		return []segment{{
			Name:       "goenv",
			Content:    segment,
			Foreground: theme.GoenvFg,
			Background: theme.GoenvBg,
		}}
	}
}
