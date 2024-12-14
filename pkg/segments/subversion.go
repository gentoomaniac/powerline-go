package segments

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

var otherModified int

func addSvnRepoStatsSegment(theme config.Theme, nChanges int, symbol string, foreground uint8, background uint8) (segments segment) {
	if nChanges > 0 {
		segments = append(segments, segment{
			Name:       "svn-status",
			Content:    fmt.Sprintf("%d%s", nChanges, symbol),
			Foreground: foreground,
			Background: background,
		})
	}
	return segments
}

func (r repoStats) SvnSegments(theme config.Theme) (segments segment) {
	segments = append(segments, addSvnRepoStatsSegment(p, r.ahead, p.symbols.RepoAhead, theme.GitAheadFg, theme.GitAheadBg)...)
	segments = append(segments, addSvnRepoStatsSegment(p, r.behind, p.symbols.RepoBehind, theme.GitBehindFg, theme.GitBehindBg)...)
	segments = append(segments, addSvnRepoStatsSegment(p, r.staged, p.symbols.RepoStaged, theme.GitStagedFg, theme.GitStagedBg)...)
	segments = append(segments, addSvnRepoStatsSegment(p, r.notStaged, p.symbols.RepoNotStaged, theme.GitNotStagedFg, theme.GitNotStagedBg)...)
	segments = append(segments, addSvnRepoStatsSegment(p, r.untracked, p.symbols.RepoUntracked, theme.GitUntrackedFg, theme.GitUntrackedBg)...)
	segments = append(segments, addSvnRepoStatsSegment(p, r.conflicted, p.symbols.RepoConflicted, theme.GitConflictedFg, theme.GitConflictedBg)...)
	segments = append(segments, addSvnRepoStatsSegment(p, r.stashed, p.symbols.RepoStashed, theme.GitStashedFg, theme.GitStashedBg)...)
	return segments
}

func runSvnCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	out, err := command.Output()
	return string(out), err
}

func parseSvnURL() (map[string]string, error) {
	info, err := runSvnCommand("svn", "info")
	if err != nil {
		return nil, errors.New("not a working copy")
	}

	svnInfo := make(map[string]string, 0)
	infos := strings.Split(info, "\n")
	if len(infos) > 1 {
		for _, line := range infos[:] {
			items := strings.Split(line, ": ")
			if len(items) >= 2 {
				svnInfo[items[0]] = items[1]
			}
		}
	}

	return svnInfo, nil
}

func ensureUnmodified(code string, stats repoStats) {
	if code != " " {
		otherModified++
	}
}

func parseSvnStatus() repoStats {
	stats := repoStats{}
	info, err := runSvnCommand("svn", "status", "-u")
	if err != nil {
		return stats
	}
	infos := strings.Split(info, "\n")
	if len(infos) > 1 {
		for _, line := range infos[:] {
			if len(line) >= 9 {
				code := line[0:1]
				switch code {
				case "?":
					stats.untracked++
				case "C":
					stats.conflicted++
				case "A", "D", "M":
					stats.notStaged++
				default:
					ensureUnmodified(code, stats)
				}
				code = line[1:2]
				switch code {
				case "C":
					stats.conflicted++
				case "M":
					stats.notStaged++
				default:
					ensureUnmodified(code, stats)
				}
				ensureUnmodified(line[2:3], stats)
				ensureUnmodified(line[3:4], stats)
				ensureUnmodified(line[4:5], stats)
				ensureUnmodified(line[5:6], stats)
				ensureUnmodified(line[6:7], stats)
				ensureUnmodified(line[7:8], stats)
				code = line[8:9]
				switch code {
				case "*":
					stats.behind++
				default:
					ensureUnmodified(code, stats)
				}
			}
		}
	}

	return stats
}

func Subversion(theme config.Theme) []segment {
	svnInfo, err := parseSvnURL()
	if err != nil {
		return []segment{}
	}

	if len(p.ignoreRepos) > 0 {
		if p.ignoreRepos[svnInfo["URL"]] || p.ignoreRepos[svnInfo["Relative URL"]] {
			return []segment{}
		}
	}

	svnStats := parseSvnStatus()

	var foreground, background uint8
	if svnStats.dirty() || otherModified > 0 {
		foreground = theme.RepoDirtyFg
		background = theme.RepoDirtyBg
	} else {
		foreground = theme.RepoCleanFg
		background = theme.RepoCleanBg
	}

	segments := segment{{
		Name:       "svn-branch",
		Content:    svnInfo["Relative URL"],
		Foreground: foreground,
		Background: background,
	}}

	segments = append(segments, svnStats.SvnSegments(p)...)
	return segments
}
