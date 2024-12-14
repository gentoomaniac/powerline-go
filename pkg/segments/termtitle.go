package segments

// Port of set_term_title segment from powerine-shell:
// https://github.com/b-ryan/powerline-shell/blob/master/powerline_shell/segments/set_term_title.py

import (
	"fmt"
	"os"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"

	pwl "github.com/gentoomaniac/powerline-go/pkg/powerline"
)

func TermTitle(theme config.Theme) []segment {
	var title string

	term := os.Getenv("TERM")
	if !(strings.Contains(term, "xterm") || strings.Contains(term, "rxvt")) {
		return []segment{}
	}

	if p.cfg.Shell == "bash" {
		title = "\\[\\e]0;\\u@\\h: \\w\\a\\]"
	} else if p.cfg.Shell == "zsh" {
		title = "%{\033]0;%n@%m: %~\007%}"
	} else {
		cwd := p.cwd
		title = fmt.Sprintf("\033]0;%s@%s: %s\007", p.username, p.hostname, cwd)
	}

	return []segment{{
		Name:           "termtitle",
		Content:        title,
		Priority:       config.MaxInteger, // do not truncate
		HideSeparators: true,              // do not draw separators
	}}
}
