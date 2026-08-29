package tui

import (
	"time"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/tui/events"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ReplLogPanel struct {
	viewport.Model
	style               lipgloss.Style
	ErrResultStyle      lipgloss.Style
	UserEntryStyle      lipgloss.Style
	OutputStyle         lipgloss.Style
	lastCommand         string
	lastCommandStarted  time.Time
	lastCommandFinished time.Time
	outputs             []string
}

func NewReplLogPanel() ReplLogPanel {
	vp := viewport.New(40, 5)
	vp.SetContent("")
	return ReplLogPanel{
		Model: vp,
		style: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 1),
	}
}

func (p ReplLogPanel) OnCommandStarted(started events.CommandStarted) {
	p.lastCommand = started.Command
	p.lastCommandStarted = started.At

}

func (p ReplLogPanel) OnCommandFinished(finished events.CommandFinished) {
	p.lastCommandFinished = finished.At
}

func (p ReplLogPanel) Update(msg tea.Msg) (ReplLogPanel, tea.Cmd) {

	switch msg.(type) {
	case tea.WindowSizeMsg:
		if m, ok := msg.(tea.WindowSizeMsg); ok {
			p.Width = m.Width - 2
			p.Height = m.Height - 10
			p.style.Height(p.Height)
			p.style.Width(p.Width)
			p.GotoBottom()
		}
	case events.CommandStarted:
		p.lastCommand = msg.(events.CommandStarted).Command
		p.lastCommandStarted = msg.(events.CommandStarted).At
	case events.CommandFinished:
		lce := msg.(events.CommandFinished)
		p.lastCommandFinished = lce.At
		p.outputs = append(p.outputs, lce.Output)
	}

	return p, nil
}

func (p ReplLogPanel) SetWidth(width int) {
	p.Width = width
}

func (p ReplLogPanel) SetHeight(height int) {
	p.Height = height
}

func (p ReplLogPanel) GotoBottom() {
	p.YOffset = p.ContentHeight() - p.Height
	if p.YOffset < 0 {
		p.YOffset = 0
	}
}

func (p ReplLogPanel) ContentHeight() int {

	return len(p.outputs)
}

func (p ReplLogPanel) View() string {
	return p.style.Render(p.Model.View())
}

func (p ReplLogPanel) AddLog(log string) {
	current := p.Model.View()
	if current != "" {
		current += "\n"
	}
	current += log
	p.SetContent(current)
	p.GotoBottom()
}
