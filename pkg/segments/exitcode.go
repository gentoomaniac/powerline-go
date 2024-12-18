package segments

import (
	"fmt"
	"strconv"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/justjanne/powerline-go/exitcode"
)

var exitCodes = map[int]string{
	// 1: generic
	1: "ERROR",
	// 2, 126, 127: common shell conventions, e.g. bash, zsh
	2:   "USAGE",
	126: "NOEXEC",
	127: "NOTFOUND",
	// 64-78: BSD, sysexits.h
	64: "USAGE",
	65: "DATAERR",
	66: "NOINPUT",
	67: "NOUSER",
	68: "NOHOST",
	69: "UNAVAILABLE",
	70: "SOFTWARE",
	71: "OSERR",
	72: "OSFILE",
	73: "CANTCREAT",
	74: "IOERR",
	75: "TEMPFAIL",
	76: "PROTOCOL",
	77: "NOPERM",
	78: "CONFIG",
}

func getMeaningFromExitCode(exitCode int) string {
	if exitCode < 128 {
		name, ok := exitCodes[exitCode]
		if ok {
			return name
		}
	} else {
		name, ok := exitcode.Signals[exitCode-128]
		if ok {
			return name
		}
	}

	return fmt.Sprintf("%d", exitCode)
}

func ExitCode(cfg config.State, align config.Alignment) []Segment {
	var meaning string
	if cfg.PrevError == 0 {
		return []Segment{}
	}
	if cfg.NumericExitCodes {
		meaning = strconv.Itoa(cfg.PrevError)
	} else {
		meaning = getMeaningFromExitCode(cfg.PrevError)
	}

	return []Segment{{
		Name:       "exit",
		Content:    meaning,
		Foreground: cfg.Theme.CmdFailedFg,
		Background: cfg.Theme.CmdFailedBg,
	}}
}
