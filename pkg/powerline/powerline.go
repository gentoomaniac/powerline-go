package powerline

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gentoomaniac/powerline-go/pkg/config"
	"github.com/gentoomaniac/powerline-go/pkg/segments"

	"github.com/mattn/go-runewidth"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
	"golang.org/x/text/width"
)

type Powerline struct {
	state          config.State
	reset          string
	priorities     map[string]int
	Segments       [][]segments.Segment
	curSegment     int
	align          config.Alignment
	rightPowerline *Powerline
}

type prioritizedSegments struct {
	i    int
	segs []segments.Segment
}

func NewPowerline(cfg config.State, align config.Alignment) *Powerline {
	p := new(Powerline)
	p.state = cfg

	if cfg.TrimAdDomain {
		usernameWithAd := strings.SplitN(cfg.Userinfo.Username, `\`, 2)
		if len(usernameWithAd) > 1 {
			// remove the Domain name from username
			cfg.Userinfo.Username = usernameWithAd[1]
		}
	}

	p.reset = fmt.Sprintf(cfg.CurrentShell().ColorTemplate, "[0m")
	p.priorities = make(map[string]int)
	for idx, priority := range cfg.Priority {
		p.priorities[priority] = len(cfg.Priority) - idx
	}
	p.align = align
	p.Segments = make([][]segments.Segment, 1)
	var mods []string
	if p.align == config.AlignLeft {
		mods = cfg.Modules
		if len(cfg.ModulesRight) > 0 {
			if p.SupportsRightModules() {
				p.rightPowerline = NewPowerline(cfg, config.AlignRight)
			} else {
				mods = append(mods, cfg.ModulesRight...)
			}
		}
	} else {
		mods = cfg.ModulesRight
	}
	initSegments(p, mods)

	return p
}

func initSegments(p *Powerline, mods []string) {
	orderedSegments := map[int][]segments.Segment{}
	c := make(chan prioritizedSegments, len(mods))
	wg := sync.WaitGroup{}
	for i, module := range mods {
		wg.Add(1)
		go func(w *sync.WaitGroup, i int, module string, c chan prioritizedSegments) {
			elem, ok := modules[module]
			if ok {
				c <- prioritizedSegments{
					i:    i,
					segs: elem(p.state, p.align),
				}
			} else {
				s, ok := segments.Plugin(p.state, module)
				if ok {
					c <- prioritizedSegments{
						i:    i,
						segs: s,
					}
				} else {
					log.Error().Str("module", module).Msg("module not found")
				}
			}
			wg.Done()
		}(&wg, i, module, c)
	}
	wg.Wait()
	close(c)
	for s := range c {
		orderedSegments[s.i] = s.segs
	}
	for i := 0; i < len(mods); i++ {
		for _, seg := range orderedSegments[i] {
			p.appendSegment(seg.Name, seg)
		}
	}
}

func (p *Powerline) color(prefix string, code uint8) string {
	if code == p.state.Theme.Reset {
		return p.reset
	}
	return fmt.Sprintf(p.state.CurrentShell().ColorTemplate, fmt.Sprintf("[%s;5;%dm", prefix, code))
}

func (p *Powerline) fgColor(code uint8) string {
	if p.state.Theme.BoldForeground {
		return p.color("1;38", code)
	} else {
		return p.color("38", code)
	}
}

func (p *Powerline) bgColor(code uint8) string {
	return p.color("48", code)
}

func (p *Powerline) appendSegment(origin string, segment segments.Segment) {
	if segment.Foreground == segment.Background && segment.Background == 0 {
		segment.Background = p.state.Theme.DefaultBg
		segment.Foreground = p.state.Theme.DefaultFg
	}
	if segment.Separator == "" {
		if p.isRightPrompt() {
			segment.Separator = p.state.Symbols().SeparatorReverse
		} else {
			segment.Separator = p.state.Symbols().Separator
		}
	}
	if segment.SeparatorForeground == 0 {
		segment.SeparatorForeground = segment.Background
	}
	segment.Priority += p.priorities[origin]
	segment.Width = segment.ComputeWidth(p.state.Condensed)
	if segment.NewLine {
		p.newRow()
	} else {
		p.Segments[p.curSegment] = append(p.Segments[p.curSegment], segment)
	}
}

func (p *Powerline) newRow() {
	if len(p.Segments[p.curSegment]) > 0 {
		p.Segments = append(p.Segments, make([]segments.Segment, 0))
		p.curSegment = p.curSegment + 1
	}
}

func (p *Powerline) truncateRow(rowNum int) {
	shellMaxLength := termWidth() * p.state.MaxWidthPercentage / 100
	row := p.Segments[rowNum]
	rowLength := 0

	if shellMaxLength > 0 {
		for _, segment := range row {
			rowLength += segment.Width
		}

		if rowLength > shellMaxLength && p.state.TruncateSegmentWidth > 0 {
			minPriorityNotTruncated := config.MaxInteger
			minPriorityNotTruncatedSegmentID := -1
			for idx, segment := range row {
				if segment.Width > p.state.TruncateSegmentWidth && segment.Priority < minPriorityNotTruncated {
					minPriorityNotTruncated = segment.Priority
					minPriorityNotTruncatedSegmentID = idx
				}
			}
			for minPriorityNotTruncatedSegmentID != -1 && rowLength > shellMaxLength {
				segment := row[minPriorityNotTruncatedSegmentID]

				rowLength -= segment.Width

				segment.Content = runewidth.Truncate(segment.Content, p.state.TruncateSegmentWidth-runewidth.StringWidth(segment.Separator)-3, "…")
				segment.Width = segment.ComputeWidth(p.state.Condensed)

				row = append(append(row[:minPriorityNotTruncatedSegmentID], segment), row[minPriorityNotTruncatedSegmentID+1:]...)
				rowLength += segment.Width

				minPriorityNotTruncated = config.MaxInteger
				minPriorityNotTruncatedSegmentID = -1
				for idx, segment := range row {
					if segment.Width > p.state.TruncateSegmentWidth && segment.Priority < minPriorityNotTruncated {
						minPriorityNotTruncated = segment.Priority
						minPriorityNotTruncatedSegmentID = idx
					}
				}
			}
		}

		for rowLength > shellMaxLength {
			minPriority := config.MaxInteger
			minPrioritySegmentID := -1
			for idx, segment := range row {
				if segment.Priority < minPriority {
					minPriority = segment.Priority
					minPrioritySegmentID = idx
				}
			}
			if minPrioritySegmentID != -1 {
				segment := row[minPrioritySegmentID]
				row = append(row[:minPrioritySegmentID], row[minPrioritySegmentID+1:]...)
				rowLength -= segment.Width
			}
		}
	}
	p.Segments[rowNum] = row
}

func (p *Powerline) numEastAsianRunes(segmentContent *string) int {
	if !p.state.EastAsianWidth {
		return 0
	}
	numEastAsianRunes := 0
	for _, r := range *segmentContent {
		switch width.LookupRune(r).Kind() {
		case width.EastAsianAmbiguous:
			numEastAsianRunes++
		case width.Neutral:
		case width.EastAsianWide:
		case width.EastAsianNarrow:
		case width.EastAsianFullwidth:
		case width.EastAsianHalfwidth:
		}
	}
	return numEastAsianRunes
}

func (p *Powerline) drawRow(rowNum int, buffer *bytes.Buffer) {
	row := p.Segments[rowNum]
	numEastAsianRunes := 0

	// Prepend padding
	if p.isRightPrompt() {
		buffer.WriteRune(' ')
	}
	for idx, segment := range row {
		if segment.HideSeparators {
			buffer.WriteString(segment.Content)
			continue
		}
		var separatorBackground string
		if p.isRightPrompt() {
			if idx == 0 {
				separatorBackground = p.reset
			} else {
				prevSegment := row[idx-1]
				separatorBackground = p.bgColor(prevSegment.Background)
			}
			buffer.WriteString(separatorBackground)
			buffer.WriteString(p.fgColor(segment.SeparatorForeground))
			buffer.WriteString(segment.Separator)
		} else {
			if idx >= len(row)-1 {
				if !p.HasRightModules() || p.SupportsRightModules() {
					separatorBackground = p.reset
				} else if p.HasRightModules() && rowNum >= len(p.Segments)-1 {
					nextSegment := p.rightPowerline.Segments[0][0]
					separatorBackground = p.bgColor(nextSegment.Background)
				}
			} else {
				nextSegment := row[idx+1]
				separatorBackground = p.bgColor(nextSegment.Background)
			}
		}
		buffer.WriteString(p.fgColor(segment.Foreground))
		buffer.WriteString(p.bgColor(segment.Background))
		if !p.state.Condensed {
			buffer.WriteRune(' ')
		}
		buffer.WriteString(segment.Content)
		numEastAsianRunes += p.numEastAsianRunes(&segment.Content)
		if !p.state.Condensed {
			buffer.WriteRune(' ')
		}
		if !p.isRightPrompt() {
			buffer.WriteString(separatorBackground)
			buffer.WriteString(p.fgColor(segment.SeparatorForeground))
			buffer.WriteString(segment.Separator)
		}
		buffer.WriteString(p.reset)
	}

	// Append padding before cursor for left-aligned prompts
	if !p.isRightPrompt() || !p.HasRightModules() {
		buffer.WriteRune(' ')
	}

	// Don't append padding for right-aligned modules
	if !p.isRightPrompt() {
		for i := 0; i < numEastAsianRunes; i++ {
			buffer.WriteRune(' ')
		}
	}
}

func (p *Powerline) Render() string {
	var buffer bytes.Buffer

	if p.state.Eval {
		if p.align == config.AlignLeft {
			buffer.WriteString(p.state.CurrentShell().EvalPromptPrefix)
		} else if p.SupportsRightModules() {
			buffer.WriteString(p.state.CurrentShell().EvalPromptRightPrefix)
		}
	}

	for rowNum := range p.Segments {
		p.truncateRow(rowNum)
		p.drawRow(rowNum, &buffer)
		if rowNum < len(p.Segments)-1 {
			// buffer.WriteRune('\n')
		}
	}

	if p.state.PromptOnNewLine {
		buffer.WriteRune('\n')

		var foreground, background uint8
		if p.state.PrevError == 0 || p.state.StaticPromptIndicator {
			foreground = p.state.Theme.CmdPassedFg
			background = p.state.Theme.CmdPassedBg
		} else {
			foreground = p.state.Theme.CmdFailedFg
			background = p.state.Theme.CmdFailedBg
		}

		buffer.WriteString(p.fgColor(foreground))
		buffer.WriteString(p.bgColor(background))
		buffer.WriteString(p.state.CurrentShell().RootIndicator)
		buffer.WriteString(p.reset)
		buffer.WriteString(p.fgColor(background))
		buffer.WriteString(p.state.Symbols().Separator)
		buffer.WriteString(p.reset)
		buffer.WriteRune(' ')
	}

	if p.state.Eval {
		switch p.align {
		case config.AlignLeft:
			buffer.WriteString(p.state.CurrentShell().EvalPromptSuffix)
			if p.SupportsRightModules() {
				buffer.WriteRune('\n')
				if !p.HasRightModules() {
					buffer.WriteString(p.state.CurrentShell().EvalPromptRightPrefix + p.state.CurrentShell().EvalPromptRightSuffix)
				}
			}
		case config.AlignRight:
			if p.SupportsRightModules() {
				buffer.Truncate(buffer.Len() - 1)
				buffer.WriteString(p.state.CurrentShell().EvalPromptRightSuffix)
			}
		}
		if p.HasRightModules() {
			buffer.WriteString(p.rightPowerline.Render())
		}
	}

	return buffer.String()
}

func (p *Powerline) HasRightModules() bool {
	return p.rightPowerline != nil && len(p.rightPowerline.Segments[0]) > 0
}

func (p *Powerline) SupportsRightModules() bool {
	return p.state.CurrentShell().EvalPromptRightPrefix != "" || p.state.CurrentShell().EvalPromptRightSuffix != ""
}

func (p *Powerline) isRightPrompt() bool {
	return p.align == config.AlignRight && p.SupportsRightModules()
}

func termWidth() int {
	termWidth, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		shellMaxLengthStr, found := os.LookupEnv("COLUMNS")
		if !found {
			return 0
		}

		shellMaxLength64, err := strconv.ParseInt(shellMaxLengthStr, 0, 64)
		if err != nil {
			return 0
		}

		termWidth = int(shellMaxLength64)
	}

	return termWidth
}
