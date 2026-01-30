package tui

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/builtins"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"strings"
)

type Model struct {
	textarea       textarea.Model
	viewport       viewport.Model
	inputStyle     lipgloss.Style
	headingStyle   lipgloss.Style
	prevInputStyle lipgloss.Style
	outputStyle    lipgloss.Style
	history        []string
	outputs        []string
	input          string
	prompt         string
	cursor         int
	err            error
	ctx            context.Context
	cancelFunc     context.CancelFunc
	runtime        *builtins.Runtime
	readFromParser io.Reader
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func InitialModel(prompt string) Model {
	ta := textarea.New()
	ta.Focus()
	ta.Prompt = prompt
	ta.CharLimit = 512
	ta.SetWidth(30)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = true
	vp := viewport.New(30, 5)
	vp.SetContent("welcome to scheme")

	//pr, pw := io.Pipe()

	ctx, cancel := context.WithCancel(context.Background())

	runtime := builtins.NewRuntime(builtins.WithOut(os.Stdout))

	// styles
	inputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	headingStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
	prevInputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	outputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

	return Model{
		viewport:       vp,
		textarea:       ta,
		history:        make([]string, 0),
		outputs:        make([]string, 0),
		prompt:         prompt,
		input:          "",
		cursor:         -1,
		inputStyle:     inputStyle,
		headingStyle:   headingStyle,
		prevInputStyle: prevInputStyle,
		outputStyle:    outputStyle,
		runtime:        runtime,

		ctx:        ctx,
		cancelFunc: cancel,
		//readFromParser: pr,
	}
}

type EvalCompleteMsg struct {
	result values.Interface
	err    error
}

type ParserIoAvailable struct {
	content  string
	canceled bool
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd                tea.Cmd
		vpCmd                tea.Cmd
		runParserCmd         tea.Cmd
		updateParseResultCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(strings.Repeat("\n", 5))
		if len(m.outputs) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.outputs, "\n")))
		}
		m.viewport.GotoBottom()

	case EvalCompleteMsg:
		return m.handleEvalComplete(msg, tiCmd, vpCmd)
	case ParserIoAvailable:
		return m.handleParserIoAvailable(msg, tiCmd, vpCmd)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			if m.cancelFunc != nil {
				m.cancelFunc()
			}
			return m, tea.Quit
		case tea.KeyCtrlR:
			// reverse search history
			return m.handleReverseHistorySearch(tiCmd, vpCmd)
		case tea.KeyEnter:
			// run parser
			return m.handleParserDispatch(msg, tiCmd, vpCmd)
		case tea.KeyUp:
			return m.handleScrollUpHistory(tiCmd, vpCmd)
		case tea.KeyDown:
			return m.handleScrollDownHistory(tiCmd, vpCmd)
		case tea.KeyBackspace:
			return m.handleBackspace(tiCmd, vpCmd)
		default:
			m.input += string(msg.Runes)
		}
	case errMsg:
		m.err = msg
	}
	if runParserCmd != nil {
		return m, tea.Batch(tiCmd, vpCmd, runParserCmd, updateParseResultCmd)
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) handleEvalComplete(msg EvalCompleteMsg, tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	now := time.Now()
	if msg.err != nil {
		m.outputs = append(m.outputs, fmt.Sprintf("%s Error: %s", now.Format(time.Kitchen), msg.err.Error()))
	}
	//if msg.result != nil {
	//	m.outputs = append(m.outputs, "Result: "+boolean.Trinary(msg.result.Type() == types.Void, "#void", msg.result.String()))
	//}
	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.outputs, "\n")))
	m.viewport.GotoBottom()
	m.input = ""
	m.cursor = -1
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) handleParserIoAvailable(msg ParserIoAvailable, tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	now := time.Now()
	if msg.content != "" {
		m.outputs = append(m.outputs, fmt.Sprintf("%s Output: %s", now.Format(time.Kitchen), msg.content))
		m.viewport.SetContent(lipgloss.NewStyle().
			Width(m.viewport.Width).
			Render(strings.Join(m.outputs, "\n")))
		m.viewport.GotoBottom()
	}
	m.input = ""
	m.cursor = -1
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) handleParserDispatch(msg tea.Msg, tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	var (
		runParserCmd, updateParseResultCmd tea.Cmd
	)
	userInput := m.input
	pr, pw := io.Pipe()
	runParserCmd = func() tea.Msg {
		m.runtime.Out = pw
		val, err := parser.EvalString(m.ctx, userInput, m.runtime)
		return EvalCompleteMsg{result: val, err: err}
	}
	updateParseResultCmd = func() tea.Msg {
		var content string
		select {
		case <-m.ctx.Done():
			_ = pr.Close()
			return ParserIoAvailable{canceled: true}
		default:
			buf := make([]byte, 1024)
			for {
				n, err := pr.Read(buf)
				if err != nil {
				}
				content = string(buf[:n])
				return ParserIoAvailable{content: content}

			}
		}
	}
	now := time.Now()
	m.history = append(m.history, m.input)
	m.outputs = append(m.outputs, fmt.Sprintf("%s Executed: %s", now.Format(time.Kitchen), m.input))
	m.viewport.
		SetContent(lipgloss.NewStyle().
			Width(m.viewport.Width).
			Render(strings.Join(m.outputs, "\n")))
	m.textarea.Reset()
	m.viewport.GotoBottom()
	return m, tea.Batch(tiCmd, vpCmd, runParserCmd, updateParseResultCmd)
}

func (m Model) handleScrollUpHistory(tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	if len(m.history) > 0 {
		m.cursor = len(m.history) - 1
		m.input = m.history[m.cursor]
		m.cursor--
	} else if m.cursor > 0 {
		m.cursor--
		m.input = m.history[m.cursor]
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) handleScrollDownHistory(tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	if len(m.history) == 0 {
		return m, tea.Batch(tiCmd, vpCmd)
	}

	if m.cursor == -1 {
		// already at the editing buffer; nothing to do
		return m, tea.Batch(tiCmd, vpCmd)
	}

	if m.cursor < len(m.history)-1 {
		// move forward in history
		m.cursor++
		m.input = m.history[m.cursor]
	} else {
		// move past the most recent entry -> restore empty input
		m.cursor = -1
		m.input = ""
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) handleBackspace(tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	if len(m.input) > 0 {
		m.input = m.input[:len(m.input)-1]
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) handleReverseHistorySearch(tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	if len(m.history) == 0 {
		return m, tea.Batch(tiCmd, vpCmd)
	}
	if m.cursor < 0 {
		m.cursor = len(m.history) - 1
	}
	for i := len(m.history) - 1; i >= 0; i-- {
		if strings.Contains(m.history[i], m.input) {
			m.cursor = i
			break
		}
	}

	m.input = m.history[m.cursor]

	return m, tea.Batch(tiCmd, vpCmd)
}

type (
	errMsg error
)

func (m Model) renderOutputs(width int) string {
	var styled []string
	for _, line := range m.outputs {
		switch {
		case strings.Contains(line, "Executed:"):
			styled = append(styled, m.prevInputStyle.Width(width).Render(line))
		case strings.Contains(line, "Output:"):
			styled = append(styled, m.outputStyle.Width(width).Render(line))
		case strings.Contains(line, "Error:"):
			// errors can be emphasized
			errStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("160"))
			styled = append(styled, errStyle.Width(width).Render(line))
		default:
			styled = append(styled, lipgloss.NewStyle().Width(width).Render(line))
		}
	}
	return strings.Join(styled, "\n")
}

func (m Model) View() string {
	var historyStr strings.Builder
	historyStr.WriteString(m.headingStyle.Render("------------- History --------") + "\n")
	for _, cmd := range m.outputs {
		// apply a simple style in View as well
		if strings.Contains(cmd, "Executed:") {
			historyStr.WriteString(m.prevInputStyle.Render(cmd) + "\n")
		} else if strings.Contains(cmd, "Output:") {
			historyStr.WriteString(m.outputStyle.Render(cmd) + "\n")
		} else {
			historyStr.WriteString(cmd + "\n")
		}
	}
	historyStr.WriteString("\n")
	historyStr.WriteString(m.headingStyle.Render("------------- Input ----------") + "\n")
	historyStr.WriteString("\n")
	historyStr.WriteString(m.inputStyle.Render(m.prompt + " " + m.input))
	return historyStr.String()
}
