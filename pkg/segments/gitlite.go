//go:build broken

package segments

import (
	"fmt"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func GitLite(theme config.Theme) []Segment {
	if len(p.ignoreRepos) > 0 {
		out, err := runGitCommand("git", "--no-optional-locks", "rev-parse", "--show-toplevel")
		if err != nil {
			return []Segment{}
		}
		out = strings.TrimSpace(out)
		if p.ignoreRepos[out] {
			return []Segment{}
		}
	}

	out, err := runGitCommand("git", "--no-optional-locks", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return []Segment{}
	}

	status := strings.TrimSpace(out)
	var branch string

	if status == "HEAD" {
		branch = getGitDetachedBranch(p)
	} else {
		branch = status
	}

	if p.cfg.GitMode != "compact" && len(p.symbols.RepoBranch) > 0 {
		branch = fmt.Sprintf("%s %s", p.symbols.RepoBranch, branch)
	}

	return []Segment{{
		Name:       "git-branch",
		Content:    branch,
		Foreground: theme.RepoCleanFg,
		Background: theme.RepoCleanBg,
	}}
}
