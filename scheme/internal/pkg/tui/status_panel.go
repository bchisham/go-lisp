package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/bchisham/collections-go/sequence"
	"github.com/bchisham/go-lisp/scheme/internal/contracts"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/tui/events"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/tui/types"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// StatusPanel is a struct that holds the state of the status panel in the TUI.
// It implements the CommandListener interface and provides methods to update the status panel based on command events and log messages.
type StatusPanel struct {
	table           table.Model
	InfoLogStyle    lipgloss.Style
	DebugLogStyle   lipgloss.Style
	WarningLogStyle lipgloss.Style

	CommandsRun        int
	Errors             int
	Warnings           int
	Info               int
	Debug              int
	LastCommandStarted time.Time
	LastCommandDone    time.Time
	LastError          error
	OpenParens         int
	LogThreshold       VerboseLevel
	logLines           []string
}

type VerboseLevel int

const (
	VerboseQuiet VerboseLevel = iota
	VerboseError VerboseLevel = iota
	VerboseWarn  VerboseLevel = iota
	VerboseInfo  VerboseLevel = iota
	VerboseDebug VerboseLevel = iota
)

func NewStatusPanel() StatusPanel {
	columns := []table.Column{
		{Title: "Commands Run", Width: 15},
		{Title: "Errors", Width: 10},
		{Title: "Warnings", Width: 10},
		{Title: "Info", Width: 10},
		{Title: "Debug", Width: 10},
		{Title: "Open Parens", Width: 15},
		{Title: "Current Command Duration", Width: 20},
		{Title: "Last Command Started", Width: 20},
	}
	rows := []table.Row{
		{"0", "0", "0", "0", "0", "0", "0s", "N/A"},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
	)
	return StatusPanel{
		table: t,
	}
}

func (s StatusPanel) OnCommandStarted(started events.CommandStarted) {
	s.LastCommandStarted = started.At
	s.Update()

}

func (s StatusPanel) OnCommandFinished(finished events.CommandFinished) {
	s.CommandsRun++
	if finished.Err != nil {
		s.Errors++
		s.LastError = finished.Err
	}
	s.LastCommandDone = finished.At
	s.Update()
}

func (s StatusPanel) Log(level VerboseLevel, format string, args ...any) StatusPanel {

	switch level {
	case VerboseError:
		s.Errors++
	case VerboseWarn:
		s.Warnings++
	case VerboseInfo:
		s.Info++
		s.Infof(format, args...)
	case VerboseDebug:
		s.Debug++
	}
	return s
}

func (s StatusPanel) Infof(format string, args ...interface{}) StatusPanel {
	s.Info++

	dspargs := s.toDisplayArgs(s.InfoLogStyle, args...)
	strArgs := sequence.MapMust(dspargs, func(d contracts.CanDisplay) string {
		return d.DisplayString()
	}).ToSlice()

	if s.shouldLog(VerboseInfo) {
		s.logLines = append(s.logLines, strings.Join(strArgs, " "))
	}
	return s
}

func (s StatusPanel) shouldLog(level VerboseLevel) bool {
	return level <= s.LogThreshold
}

func (s StatusPanel) Update() {
	r := table.Row{
		fmt.Sprintf("%03d", s.CommandsRun),
		fmt.Sprintf("%03d", s.Errors),
		fmt.Sprintf("%03d", s.Warnings),
		fmt.Sprintf("%03d", s.Info),
		fmt.Sprintf("%03d", s.Debug),
		fmt.Sprintf("%03d", s.OpenParens),
		s.LastCommandDone.Sub(s.LastCommandStarted).String(),
		s.LastCommandStarted.Format(time.RFC822),
	}
	s.table.SetRows([]table.Row{r})
}

func (s StatusPanel) View() string {
	return s.table.View()
}

func (s StatusPanel) toDisplayArgs(style lipgloss.Style, args ...interface{}) []contracts.CanDisplay {
	iArgs := sequence.FromSlice(args)
	xform := sequence.NewTransformer[interface{}, contracts.CanDisplay](iArgs)
	return xform.TransformMust(types.NewAnyCanDisplayFunc[interface{}](
		types.WithDisplayStyle(style),
	)).ToSlice()
}
