package segments

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

var otherModified int

func addSvnRepoStatsSegment(nChanges int, symbol string, foreground uint8, background uint8) (segments []Segment) {
	if nChanges > 0 {
		segments = append(segments, Segment{
			Name:       "svn-status",
			Content:    fmt.Sprintf("%d%s", nChanges, symbol),
			Foreground: foreground,
			Background: background,
		})
	}
	return segments
}

func (r repoStats) SvnSegments(cfg config.State) (segments []Segment) {
	segments = append(segments, addSvnRepoStatsSegment(r.ahead, cfg.Symbols().RepoAhead, cfg.Theme.GitAheadFg, cfg.Theme.GitAheadBg)...)
	segments = append(segments, addSvnRepoStatsSegment(r.behind, cfg.Symbols().RepoBehind, cfg.Theme.GitBehindFg, cfg.Theme.GitBehindBg)...)
	segments = append(segments, addSvnRepoStatsSegment(r.staged, cfg.Symbols().RepoStaged, cfg.Theme.GitStagedFg, cfg.Theme.GitStagedBg)...)
	segments = append(segments, addSvnRepoStatsSegment(r.notStaged, cfg.Symbols().RepoNotStaged, cfg.Theme.GitNotStagedFg, cfg.Theme.GitNotStagedBg)...)
	segments = append(segments, addSvnRepoStatsSegment(r.untracked, cfg.Symbols().RepoUntracked, cfg.Theme.GitUntrackedFg, cfg.Theme.GitUntrackedBg)...)
	segments = append(segments, addSvnRepoStatsSegment(r.conflicted, cfg.Symbols().RepoConflicted, cfg.Theme.GitConflictedFg, cfg.Theme.GitConflictedBg)...)
	segments = append(segments, addSvnRepoStatsSegment(r.stashed, cfg.Symbols().RepoStashed, cfg.Theme.GitStashedFg, cfg.Theme.GitStashedBg)...)
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

func Subversion(cfg config.State, align config.Alignment) []Segment {
	svnInfo, err := parseSvnURL()
	if err != nil {
		return []Segment{}
	}

	for _, repo := range cfg.IgnoreRepos {
		if svnInfo["URL"] == repo || svnInfo["Relative URL"] == repo {
			return []Segment{}
		}
	}

	svnStats := parseSvnStatus()

	var foreground, background uint8
	if svnStats.dirty() || otherModified > 0 {
		foreground = cfg.Theme.RepoDirtyFg
		background = cfg.Theme.RepoDirtyBg
	} else {
		foreground = cfg.Theme.RepoCleanFg
		background = cfg.Theme.RepoCleanBg
	}

	segments := []Segment{{
		Name:       "svn-branch",
		Content:    svnInfo["Relative URL"],
		Foreground: foreground,
		Background: background,
	}}

	segments = append(segments, svnStats.SvnSegments(cfg)...)
	return segments
}
