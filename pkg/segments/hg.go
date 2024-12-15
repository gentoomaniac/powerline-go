package segments

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func getHgStatus() (bool, bool, bool) {
	hasModifiedFiles := false
	hasUntrackedFiles := false
	hasMissingFiles := false

	out, err := exec.Command("hg", "status").Output()
	if err == nil {
		output := strings.Split(string(out), "\n")
		for _, line := range output {
			if line != "" {
				if line[0] == '?' {
					hasUntrackedFiles = true
				} else if line[0] == '!' {
					hasMissingFiles = true
				} else {
					hasModifiedFiles = true
				}
			}
		}
	}
	return hasModifiedFiles, hasUntrackedFiles, hasMissingFiles
}

func Hg(cfg config.Config, align config.Alignment) []Segment {
	out, _ := exec.Command("hg", "branch").Output()
	output := strings.SplitN(string(out), "\n", 2)
	if !(len(output) > 0 && output[0] != "") {
		return []Segment{}
	}
	branch := output[0]
	hasModifiedFiles, hasUntrackedFiles, hasMissingFiles := getHgStatus()

	var foreground, background uint8
	var content string
	if hasModifiedFiles || hasUntrackedFiles || hasMissingFiles {
		foreground = cfg.Theme.RepoDirtyFg
		background = cfg.Theme.RepoDirtyBg

		extra := ""

		if hasUntrackedFiles {
			extra += "+"
		}

		if hasMissingFiles {
			extra += "!"
		}

		content = fmt.Sprintf("%s %s", branch, extra)
	} else {
		foreground = cfg.Theme.RepoCleanFg
		background = cfg.Theme.RepoCleanBg

		content = branch
	}

	return []Segment{{
		Name:       "hg",
		Content:    content,
		Foreground: foreground,
		Background: background,
	}}
}
