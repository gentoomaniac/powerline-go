package segments

import (
	"os"
	"sort"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

const ellipsis = "\u2026"

type pathSegment struct {
	path     string
	home     bool
	root     bool
	ellipsis bool
	alias    bool
}

type byRevLength []string

func (s byRevLength) Len() int {
	return len(s)
}

func (s byRevLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byRevLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func maybeAliasPathSegments(cfg config.Config, pathSegments []pathSegment) []pathSegment {
	pathSeparator := string(os.PathSeparator)

	if cfg.PathAliases == nil || len(cfg.PathAliases) == 0 {
		return pathSegments
	}

	keys := make([]string, len(cfg.PathAliases))
	for k := range cfg.PathAliases {
		keys = append(keys, k)
	}
	sort.Sort(byRevLength(keys))

Aliases:
	for _, k := range keys {
		// This turns a string like "foo/bar/baz" into an array of strings.
		path := strings.Split(strings.Trim(k, pathSeparator), pathSeparator)

		// If the path has 3 elements, we know we should look at pathSegments
		// in 3-element chunks.
		size := len(path)
		// If there aren't that many segments in our path we can skip to the
		// next alias.
		if size > len(pathSegments) {
			continue Aliases
		}

		alias := cfg.PathAliases[k]

	Segments:
		// We want to see if that array of strings exists in pathSegments.
		for i := range pathSegments {
			// This is the upper index that we would look at. So if i is 0,
			// then we'd look at pathSegments[0,1,2], then [1,2,3], etc.. If i
			// is 2, we'd look at pathSegments[2,3,4] and so on.
			max := (i + size) - 1

			// But if the upper index is out of bounds we can short-circuit
			// and move on to the next alias.
			if max > (len(pathSegments)-i)-1 {
				continue Aliases
			}

			// Then we loop over the indices in path and compare the
			// elements. If any element doesn't match we can move on to the
			// next index in pathSegments.
			for j := range path {
				if path[j] != pathSegments[i+j].path {
					continue Segments
				}
			}

			// They all matched! That means we can replace this slice with our
			// alias and skip to the next alias.
			pathSegments = append(
				pathSegments[:i],
				append(
					[]pathSegment{{
						path:  alias,
						alias: true,
					}},
					pathSegments[max+1:]...,
				)...,
			)
			continue Aliases
		}
	}

	return pathSegments
}

func cwdToPathSegments(cfg config.Config, cwd string) []pathSegment {
	pathSeparator := string(os.PathSeparator)
	pathSegments := make([]pathSegment, 0)

	if cwd == cfg.Userinfo.HomeDir {
		pathSegments = append(pathSegments, pathSegment{
			path: "~",
			home: true,
		})
		cwd = ""
	} else if strings.HasPrefix(cwd, cfg.Userinfo.HomeDir+pathSeparator) {
		pathSegments = append(pathSegments, pathSegment{
			path: "~",
			home: true,
		})
		cwd = cwd[len(cfg.Userinfo.HomeDir):]
	} else if cwd == pathSeparator {
		pathSegments = append(pathSegments, pathSegment{
			path: pathSeparator,
			root: true,
		})
	}

	cwd = strings.Trim(cwd, pathSeparator)
	names := strings.Split(cwd, pathSeparator)
	if names[0] == "" {
		names = names[1:]
	}

	for _, name := range names {
		pathSegments = append(pathSegments, pathSegment{
			path: name,
		})
	}

	return maybeAliasPathSegments(cfg, pathSegments)
}

func maybeShortenName(cfg config.Config, pathSegment string) string {
	if cfg.CwdMaxDirSize > 0 && len(pathSegment) > cfg.CwdMaxDirSize {
		return pathSegment[:cfg.CwdMaxDirSize]
	}
	return pathSegment
}

func getColor(cfg config.Config, pathSegment pathSegment, isLastDir bool) (uint8, uint8, bool) {
	if pathSegment.home && cfg.Theme.HomeSpecialDisplay {
		return cfg.Theme.HomeFg, cfg.Theme.HomeBg, true
	} else if pathSegment.alias {
		return cfg.Theme.AliasFg, cfg.Theme.AliasBg, true
	} else if isLastDir {
		return cfg.Theme.CwdFg, cfg.Theme.PathBg, false
	}
	return cfg.Theme.PathFg, cfg.Theme.PathBg, false
}

func Cwd(cfg config.Config, align config.Alignment) (segments []Segment) {
	cwd := cfg.Cwd

	switch cfg.CwdMode {
	case "plain":
		if strings.HasPrefix(cwd, cfg.Userinfo.HomeDir) {
			cwd = "~" + cwd[len(cfg.Userinfo.HomeDir):]
		}

		segments = append(segments, Segment{
			Name:       "cwd",
			Content:    escapeVariables(cfg, cwd),
			Foreground: cfg.Theme.CwdFg,
			Background: cfg.Theme.PathBg,
		})
	default:
		pathSegments := cwdToPathSegments(cfg, cwd)

		if cfg.CwdMode == "dironly" {
			pathSegments = pathSegments[len(pathSegments)-1:]
		} else {
			maxDepth := cfg.CwdMaxDepth
			if maxDepth <= 0 {
				log.Warn().Msg("Ignoring -cwd-max-depth argument since it's smaller than or equal to 0")
			} else if len(pathSegments) > maxDepth {
				var nBefore int
				if maxDepth > 2 {
					nBefore = 2
				} else {
					nBefore = maxDepth - 1
				}
				firstPart := pathSegments[:nBefore]
				secondPart := pathSegments[len(pathSegments)+nBefore-maxDepth:]

				pathSegments = make([]pathSegment, 0)
				pathSegments = append(pathSegments, firstPart...)
				pathSegments = append(pathSegments, pathSegment{
					path:     ellipsis,
					ellipsis: true,
				})
				pathSegments = append(pathSegments, secondPart...)
			}

			if cfg.CwdMode == "semifancy" && len(pathSegments) > 1 {
				var path string
				for idx, pathSegment := range pathSegments {
					if pathSegment.home || pathSegment.alias {
						continue
					}
					path += pathSegment.path
					if idx != len(pathSegments)-1 {
						path += string(os.PathSeparator)
					}
				}
				first := pathSegments[0]
				pathSegments = make([]pathSegment, 0)
				if first.home || first.alias {
					pathSegments = append(pathSegments, first)
				}
				pathSegments = append(pathSegments, pathSegment{
					path: path,
				})
			}
		}

		for idx, pathSegment := range pathSegments {
			isLastDir := idx == len(pathSegments)-1
			foreground, background, special := getColor(cfg, pathSegment, isLastDir)

			segment := Segment{
				Content:    escapeVariables(cfg, maybeShortenName(cfg, pathSegment.path)),
				Foreground: foreground,
				Background: background,
			}

			if !special {
				// TODO: the supports check needs to get done in powerline
				// && p.supportsRightModules()
				if align == config.AlignRight && idx != 0 {
					segment.Separator = cfg.Symbols().SeparatorReverseThin
					segment.SeparatorForeground = cfg.Theme.SeparatorFg
				} else if (align == config.AlignLeft) && !isLastDir {
					segment.Separator = cfg.Symbols().SeparatorThin
					segment.SeparatorForeground = cfg.Theme.SeparatorFg
				}
			}

			segment.Name = "cwd-path"
			if isLastDir {
				segment.Name = "cwd"
			}

			segments = append(segments, segment)
		}
	}
	return segments
}
