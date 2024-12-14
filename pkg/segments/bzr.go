package segments

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func getBzrStatus() (bool, bool, bool) {
	hasModifiedFiles := false
	hasUntrackedFiles := false
	hasMissingFiles := false

	out, err := exec.Command("bzr", "status").Output()
	if err == nil {
		output := strings.Split(string(out), "\n")
		for _, line := range output {
			if line != "" {
				if line == "unknown:" {
					hasUntrackedFiles = true
				} else if line == "missing:" {
					hasMissingFiles = true
				} else {
					hasModifiedFiles = true
				}
			}
		}
	}
	return hasModifiedFiles, hasUntrackedFiles, hasMissingFiles
}

func Bzr(cfg config.Config) []Segment {
	out, _ := exec.Command("bzr", "nick").Output()
	output := strings.SplitN(string(out), "\n", 2)
	if len(output) > 0 && output[0] != "" {
		branch := output[0]
		hasModifiedFiles, hasUntrackedFiles, hasMissingFiles := getBzrStatus()

		var foreground, background uint8
		var content string
		if hasModifiedFiles || hasUntrackedFiles || hasMissingFiles {
			foreground = cfg.SelectedTheme().RepoDirtyFg
			background = cfg.SelectedTheme().RepoDirtyBg

			extra := ""

			if hasUntrackedFiles {
				extra += "+"
			}

			if hasMissingFiles {
				extra += "!"
			}

			if hasUntrackedFiles {
				extra += "?"
			}

			content = fmt.Sprintf("%s %s", branch, extra)
		} else {
			foreground = cfg.SelectedTheme().RepoCleanFg
			background = cfg.SelectedTheme().RepoCleanBg

			content = branch
		}

		return []Segment{{
			Name:       "bzr",
			Content:    content,
			Foreground: foreground,
			Background: background,
		}}
	}
	return []Segment{}
}
