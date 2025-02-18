package segments

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func runRvmCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	out, err := command.Output()
	return string(out), err
}

// check RUBY_VERSION variable
func checkEnvForRubyVersion() (string, error) {
	rubyVersion := os.Getenv("RUBY_VERSION")
	if len(rubyVersion) <= 0 {
		return "", errors.New("not found in RUBY_VERSION")
	}
	return rubyVersion, nil
}

// check GEM_HOME variable for gemset information
func checkEnvForRubyGemset() (string, error) {
	gemHomeSegments := strings.Split(os.Getenv("GEM_HOME"), "@")

	if len(gemHomeSegments) <= 1 {
		return "", errors.New("gemset not found in GEM_HOME")
	}

	return gemHomeSegments[1], nil
}

// retrieve ruby version from RVM
func checkForRvmOutput() (string, error) {
	// ask RVM what the current ruby version is
	out, err := runRvmCommand("rvm", "current")
	if err != nil {
		return "", errors.New("not found in RVM output")
	}
	items := strings.Split(out, " ")
	if len(items) <= 0 {
		return "", errors.New("not found in RVM output")
	}

	return items[0], nil
}

func Rvm(cfg config.State, align config.Alignment) []Segment {
	var (
		segment string
		err     error
	)

	segment, err = checkEnvForRubyVersion()
	if err != nil {
		segment, err = checkForRubyVersionFileInTree()
	}
	if err != nil {
		segment, err = checkForRvmOutput()
	}
	if err != nil {
		return []Segment{}
	}

	// Remove explicit "ruby-" prefix from segment because it's superfluous
	segment_components := strings.Split(segment, "-")
	if len(segment_components) > 1 {
		segment = segment_components[1]
	}

	// If gemset is missing from segment, get that info from the environment
	segment_components = strings.Split(segment, "@")
	if len(segment_components) < 2 {
		gemset, err := checkEnvForRubyGemset()
		if err == nil && gemset != "" {
			segment = segment + "@" + gemset
		}
	}

	return []Segment{{
		Name:       "rvm",
		Content:    cfg.Symbols().RvmIndicator + " " + segment,
		Foreground: cfg.Theme.RvmFg,
		Background: cfg.Theme.RvmBg,
	}}
}
