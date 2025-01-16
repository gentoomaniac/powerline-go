package segments

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/gentoomaniac/powerline-go/pkg/signals"
	"github.com/rs/zerolog/log"
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
		name := signals.ResolveSignalCode(exitCode - 128)
		if name != "" {
			return name
		}

	}

	return fmt.Sprintf("%d", exitCode)
}

func getExitCode() (int, error) {
	val := os.Getenv("?")
	code, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return int(code), nil
	}
	val = os.Getenv("status")
	code, err = strconv.ParseInt(val, 10, 32)
	if err != nil {
		return int(code), nil
	}
	return -255, fmt.Errorf("no exit code found")
}

func ExitCode(cfg config.State, align config.Alignment) []Segment {
	var meaning string
	code, err := getExitCode()
	log.Debug().Msgf("%d", code)
	if code == 0 || err != nil {
		return []Segment{}
	}
	if cfg.NumericExitCodes {
		meaning = strconv.Itoa(code)
	} else {
		meaning = getMeaningFromExitCode(code)
	}

	return []Segment{{
		Name:       "exit",
		Content:    meaning,
		Foreground: cfg.Theme.CmdFailedFg,
		Background: cfg.Theme.CmdFailedBg,
	}}
}
