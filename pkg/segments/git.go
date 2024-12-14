package segments

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

func homeEnvName() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	}
	return env
}

type repoStats struct {
	ahead      int
	behind     int
	untracked  int
	notStaged  int
	staged     int
	conflicted int
	stashed    int
}

func (r repoStats) dirty() bool {
	return r.untracked+r.notStaged+r.staged+r.conflicted > 0
}

func (r repoStats) any() bool {
	return r.ahead+r.behind+r.untracked+r.notStaged+r.staged+r.conflicted+r.stashed > 0
}

func addRepoStatsSegment(nChanges int, symbol string, foreground uint8, background uint8) []Segment {
	if nChanges > 0 {
		return []Segment{{
			Name:       "git-status",
			Content:    fmt.Sprintf("%d%s", nChanges, symbol),
			Foreground: foreground,
			Background: background,
		}}
	}
	return []Segment{}
}

func (r repoStats) GitSegments(cfg config.Config) (segments []Segment) {
	segments = append(segments, addRepoStatsSegment(r.ahead, cfg.Symbols().RepoAhead, cfg.SelectedTheme().GitAheadFg, cfg.SelectedTheme().GitAheadBg)...)
	segments = append(segments, addRepoStatsSegment(r.behind, cfg.Symbols().RepoBehind, cfg.SelectedTheme().GitBehindFg, cfg.SelectedTheme().GitBehindBg)...)
	segments = append(segments, addRepoStatsSegment(r.staged, cfg.Symbols().RepoStaged, cfg.SelectedTheme().GitStagedFg, cfg.SelectedTheme().GitStagedBg)...)
	segments = append(segments, addRepoStatsSegment(r.notStaged, cfg.Symbols().RepoNotStaged, cfg.SelectedTheme().GitNotStagedFg, cfg.SelectedTheme().GitNotStagedBg)...)
	segments = append(segments, addRepoStatsSegment(r.untracked, cfg.Symbols().RepoUntracked, cfg.SelectedTheme().GitUntrackedFg, cfg.SelectedTheme().GitUntrackedBg)...)
	segments = append(segments, addRepoStatsSegment(r.conflicted, cfg.Symbols().RepoConflicted, cfg.SelectedTheme().GitConflictedFg, cfg.SelectedTheme().GitConflictedBg)...)
	segments = append(segments, addRepoStatsSegment(r.stashed, cfg.Symbols().RepoStashed, cfg.SelectedTheme().GitStashedFg, cfg.SelectedTheme().GitStashedBg)...)
	return
}

func addRepoStatsSymbol(nChanges int, symbol string, GitMode string) string {
	if nChanges > 0 {
		if GitMode == "simple" {
			return symbol
		} else if GitMode == "compact" {
			return fmt.Sprintf(" %d%s", nChanges, symbol)
		} else {
			return symbol
		}
	}
	return ""
}

func (r repoStats) GitSymbols(cfg config.Config) string {
	var info string
	info += addRepoStatsSymbol(r.ahead, cfg.Symbols().RepoAhead, cfg.GitMode)
	info += addRepoStatsSymbol(r.behind, cfg.Symbols().RepoBehind, cfg.GitMode)
	info += addRepoStatsSymbol(r.staged, cfg.Symbols().RepoStaged, cfg.GitMode)
	info += addRepoStatsSymbol(r.notStaged, cfg.Symbols().RepoNotStaged, cfg.GitMode)
	info += addRepoStatsSymbol(r.untracked, cfg.Symbols().RepoUntracked, cfg.GitMode)
	info += addRepoStatsSymbol(r.conflicted, cfg.Symbols().RepoConflicted, cfg.GitMode)
	info += addRepoStatsSymbol(r.stashed, cfg.Symbols().RepoStashed, cfg.GitMode)
	return info
}

var branchRegex = regexp.MustCompile(`^## (?P<local>\S+?)(\.{3}(?P<remote>\S+?)( \[(ahead (?P<ahead>\d+)(, )?)?(behind (?P<behind>\d+))?])?)?$`)

func groupDict(pattern *regexp.Regexp, haystack string) map[string]string {
	match := pattern.FindStringSubmatch(haystack)
	result := make(map[string]string)
	if len(match) > 0 {
		for i, name := range pattern.SubexpNames() {
			if i != 0 {
				result[name] = match[i]
			}
		}
	}
	return result
}

var gitProcessEnv = func() []string {
	homeEnv := homeEnvName()
	home, _ := os.LookupEnv(homeEnv)
	path, _ := os.LookupEnv("PATH")
	env := map[string]string{
		"LANG":  "C",
		homeEnv: home,
		"PATH":  path,
	}
	result := make([]string, 0)
	for key, value := range env {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}()

func runGitCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	command.Env = gitProcessEnv
	out, err := command.Output()
	return string(out), err
}

func parseGitBranchInfo(status []string) map[string]string {
	return groupDict(branchRegex, status[0])
}

func getGitDetachedBranch(cfg config.Config) string {
	out, err := runGitCommand("git", "--no-optional-locks", "rev-parse", "--short", "HEAD")
	if err != nil {
		out, err := runGitCommand("git", "--no-optional-locks", "symbolic-ref", "--short", "HEAD")
		if err != nil {
			return "Error"
		}
		return strings.SplitN(out, "\n", 2)[0]
	}
	detachedRef := strings.SplitN(out, "\n", 2)
	return fmt.Sprintf("%s %s", cfg.Symbols().RepoDetached, detachedRef[0])
}

func parseGitStats(status []string) repoStats {
	stats := repoStats{}
	if len(status) > 1 {
		for _, line := range status[1:] {
			if len(line) > 2 {
				code := line[:2]
				switch code {
				case "??":
					stats.untracked++
				case "DD", "AU", "UD", "UA", "DU", "AA", "UU":
					stats.conflicted++
				default:
					if code[0] != ' ' {
						stats.staged++
					}

					if code[1] != ' ' {
						stats.notStaged++
					}
				}
			}
		}
	}
	return stats
}

func repoRoot(path string) (string, error) {
	out, err := runGitCommand("git", "--no-optional-locks", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func indexSize(root string) (int64, error) {
	fileInfo, err := os.Stat(path.Join(root, ".git", "index"))
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

func Git(cfg config.Config, align config.Alignment) []Segment {
	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("could not determine current working directory")
	}
	repoRoot, err := repoRoot(cwd)
	if err != nil {
		return []Segment{}
	}

	for _, repo := range cfg.IgnoreRepos {
		if repoRoot == repo {
			return []Segment{}
		}
	}

	args := []string{
		"--no-optional-locks", "status", "--porcelain", "-b", "--ignore-submodules",
	}

	if cfg.GitAssumeUnchangedSize > 0 {
		indexSize, _ := indexSize(cwd)
		if indexSize > (cfg.GitAssumeUnchangedSize * 1024) {
			args = append(args, "-uno")
		}
	}

	out, err := runGitCommand("git", args...)
	if err != nil {
		return []Segment{}
	}

	status := strings.Split(out, "\n")
	stats := parseGitStats(status)
	branchInfo := parseGitBranchInfo(status)
	var branch string

	if branchInfo["local"] != "" {
		ahead, _ := strconv.ParseInt(branchInfo["ahead"], 10, 32)
		stats.ahead = int(ahead)

		behind, _ := strconv.ParseInt(branchInfo["behind"], 10, 32)
		stats.behind = int(behind)

		branch = branchInfo["local"]
	} else {
		branch = getGitDetachedBranch(cfg)
	}

	if len(cfg.Symbols().RepoBranch) > 0 {
		branch = fmt.Sprintf("%s %s", cfg.Symbols().RepoBranch, branch)
	}

	var foreground, background uint8
	if stats.dirty() {
		foreground = cfg.SelectedTheme().RepoDirtyFg
		background = cfg.SelectedTheme().RepoDirtyBg
	} else {
		foreground = cfg.SelectedTheme().RepoCleanFg
		background = cfg.SelectedTheme().RepoCleanBg
	}

	segments := []Segment{{
		Name:       "git-branch",
		Content:    branch,
		Foreground: foreground,
		Background: background,
	}}

	stashEnabled := true
	for _, stat := range cfg.GitDisableStats {
		// "ahead, behind, staged, notStaged, untracked, conflicted, stashed"
		switch stat {
		case "ahead":
			stats.ahead = 0
		case "behind":
			stats.behind = 0
		case "staged":
			stats.staged = 0
		case "notStaged":
			stats.notStaged = 0
		case "untracked":
			stats.untracked = 0
		case "conflicted":
			stats.conflicted = 0
		case "stashed":
			stats.stashed = 0
			stashEnabled = false
		}
	}

	if stashEnabled {
		out, err = runGitCommand("git", "--no-optional-locks", "rev-list", "-g", "refs/stash")
		if err == nil {
			stats.stashed = strings.Count(out, "\n")
		}
	}

	if cfg.GitMode == "simple" {
		if stats.any() {
			segments[0].Content += " " + stats.GitSymbols(cfg)
		}
	} else if cfg.GitMode == "compact" {
		if stats.any() {
			segments[0].Content += stats.GitSymbols(cfg)
		}
	} else { // fancy
		segments = append(segments, stats.GitSegments(cfg)...)
	}

	return segments
}
