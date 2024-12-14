//go:build broken

package segments

// Port of set_term_title segment from powerine-shell:
// https://github.com/b-ryan/powerline-shell/blob/master/powerline_shell/segments/set_term_title.py

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/rs/zerolog/log"
)

func TermTitle(cfg config.Config, align config.Alignment) []Segment {
	var title string
	user, err := user.Current()
	if err != nil {
		log.Error().Err(err).Msg("failed getting userinfo")
		return []Segment{}
	}

	term := os.Getenv("TERM")
	if !(strings.Contains(term, "xterm") || strings.Contains(term, "rxvt")) {
		return []Segment{}
	}

	if p.cfg.Shell == "bash" {
		title = "\\[\\e]0;\\u@\\h: \\w\\a\\]"
	} else if p.cfg.Shell == "zsh" {
		title = "%{\033]0;%n@%m: %~\007%}"
	} else {
		cwd := p.cwd
		title = fmt.Sprintf("\033]0;%s@%s: %s\007", user.Username, p.hostname, cwd)
	}

	// TODO: This doesn't actually show the term title but an arbitrary string

	return []Segment{{
		Name:           "termtitle",
		Content:        title,
		Priority:       config.MaxInteger, // do not truncate
		HideSeparators: true,              // do not draw separators
	}}
}
