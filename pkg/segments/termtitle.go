package segments

// Port of set_term_title segment from powerine-shell:
// https://github.com/b-ryan/powerline-shell/blob/master/powerline_shell/segments/set_term_title.py

import (
	"fmt"
	"os"
	"strings"

	"github.com/gentoomaniac/powerline-go/pkg/config"
)

func TermTitle(cfg config.Config, align config.Alignment) []Segment {
	var title string

	term := os.Getenv("TERM")
	if !(strings.Contains(term, "xterm") || strings.Contains(term, "rxvt")) {
		return []Segment{}
	}

	if cfg.Shell == "bash" {
		title = "\\[\\e]0;\\u@\\h\\]"
	} else if cfg.Shell == "zsh" {
		title = "%{\033]0;%n@%m: %~\007%}"
	} else {
		title = fmt.Sprintf("\033]0;%s@%s: %s\007", cfg.Userinfo.Username, cfg.Hostname, cfg.Cwd)
	}

	return []Segment{{
		Name:           "termtitle",
		Content:        title,
		Priority:       config.MaxInteger, // do not truncate
		HideSeparators: true,              // do not draw separators
	}}
}
