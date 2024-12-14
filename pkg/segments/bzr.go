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

func Bzr(theme config.Theme) []segment {
	out, _ := exec.Command("bzr", "nick").Output()
	output := strings.SplitN(string(out), "\n", 2)
	if len(output) > 0 && output[0] != "" {
		branch := output[0]
		hasModifiedFiles, hasUntrackedFiles, hasMissingFiles := getBzrStatus()

		var foreground, background uint8
		var content string
		if hasModifiedFiles || hasUntrackedFiles || hasMissingFiles {
			foreground = theme.RepoDirtyFg
			background = theme.RepoDirtyBg

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
			foreground = theme.RepoCleanFg
			background = theme.RepoCleanBg

			content = branch
		}

		return []segment{{
			Name:       "bzr",
			Content:    content,
			Foreground: foreground,
			Background: background,
		}}
	}
	return []segment{}
}
