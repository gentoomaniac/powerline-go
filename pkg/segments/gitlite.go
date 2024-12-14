package segments

import (
	"fmt"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func GitLite(theme config.Theme) []segment {
	if len(p.ignoreRepos) > 0 {
		out, err := runGitCommand("git", "--no-optional-locks", "rev-parse", "--show-toplevel")
		if err != nil {
			return []segment{}
		}
		out = strings.TrimSpace(out)
		if p.ignoreRepos[out] {
			return []segment{}
		}
	}

	out, err := runGitCommand("git", "--no-optional-locks", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return []segment{}
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

	return []segment{{
		Name:       "git-branch",
		Content:    branch,
		Foreground: theme.RepoCleanFg,
		Background: theme.RepoCleanBg,
	}}
}
