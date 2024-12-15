package segments

import (
	"fmt"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func GitLite(cfg config.Config, align config.Alignment) []Segment {
	if len(cfg.IgnoreRepos) > 0 {
		out, err := runGitCommand("git", "--no-optional-locks", "rev-parse", "--show-toplevel")
		if err != nil {
			return []Segment{}
		}
		out = strings.TrimSpace(out)
		for _, repo := range cfg.IgnoreRepos {
			if out == repo {
				return []Segment{}
			}
		}
	}

	out, err := runGitCommand("git", "--no-optional-locks", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return []Segment{}
	}

	status := strings.TrimSpace(out)
	var branch string

	if status == "HEAD" {
		branch = getGitDetachedBranch(cfg)
	} else {
		branch = status
	}

	if cfg.GitMode != "compact" && len(cfg.Symbols().RepoBranch) > 0 {
		branch = fmt.Sprintf("%s %s", cfg.Symbols().RepoBranch, branch)
	}

	return []Segment{{
		Name:       "git-branch",
		Content:    branch,
		Foreground: cfg.Theme.RepoCleanFg,
		Background: cfg.Theme.RepoCleanBg,
	}}
}
