package tui

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/builtins"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textarea           textarea.Model
	viewport           viewport.Model
	cursorStyle        lipgloss.Style
	cursorLineStyle    lipgloss.Style
	inputStyle         lipgloss.Style
	headingStyle       lipgloss.Style
	prevInputStyle     lipgloss.Style
	outputStyle        lipgloss.Style
	focusedBorderStyle lipgloss.Style
	blurredBorderStyle lipgloss.Style
	history            []string
	outputs            []string
	input              string
	pendingCommand     strings.Builder
	prompt             string
	cursor             int
	err                error
	ctx                context.Context
	cancelFunc         context.CancelFunc
	runtime            *builtins.Runtime
	readFromParser     io.Reader
	openParens         int
	keymap             keymap
	help               help.Model
}

func InitialModel(prompt string) Model {
	//init styles
	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("32"))
	cursorLineStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("240"))

	focusedBorderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
	blurredBorderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	headingStyle := lipgloss.NewStyle().
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("69"))
	prevInputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	outputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250"))

	ta := textarea.New()
	ta.Focus()

	ta.SetPromptFunc(len(prompt), func(lineIdx int) string {
		if lineIdx == 0 {
			return prompt
		}
		return ""
	})
	ta.Placeholder = "Enter a scheme expression..."
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color("235"))
	ta.BlurredStyle.Base = blurredBorderStyle
	ta.FocusedStyle.Base = focusedBorderStyle
	ta.CharLimit = 1024
	ta.SetWidth(80)
	ta.SetHeight(10)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = true
	vp := viewport.New(30, 5)
	vp.SetContent("welcome to scheme")

	ctx, cancel := context.WithCancel(context.Background())

	runtime := builtins.NewRuntime(
		builtins.WithOut(os.Stdout),
		builtins.WithEvaluatorCallback(parser.DefaultExpressionEvaluator()))

	return Model{
		viewport:           vp,
		textarea:           ta,
		history:            make([]string, 0),
		outputs:            make([]string, 0),
		prompt:             prompt,
		input:              "",
		cursor:             -1,
		inputStyle:         inputStyle,
		headingStyle:       headingStyle,
		prevInputStyle:     prevInputStyle,
		outputStyle:        outputStyle,
		runtime:            runtime,
		openParens:         0,
		ctx:                ctx,
		cancelFunc:         cancel,
		cursorStyle:        cursorStyle,
		cursorLineStyle:    cursorLineStyle,
		focusedBorderStyle: focusedBorderStyle,
		blurredBorderStyle: blurredBorderStyle,
		keymap: keymap{
			Up: key.NewBinding(
				key.WithKeys("up"),
				key.WithHelp("↑", "previous command"),
			),
			Down: key.NewBinding(
				key.WithKeys("down"),
				key.WithHelp("↓", "next command"),
			),
			Enter: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "pend command"),
			),
			ReverseSearch: key.NewBinding(
				key.WithKeys("ctrl+r"),
				key.WithHelp("ctrl+r", "reverse search history"),
			),
			ScrollUpHistory: key.NewBinding(
				key.WithKeys("ctrl+up"),
				key.WithHelp("ctrl+↑", "scroll up history"),
			),
			ScrollDownHistory: key.NewBinding(
				key.WithKeys("ctrl+down"),
				key.WithHelp("ctrl+↓", "scroll down history"),
			),
			Submit: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "submit command"),
			),
			Exit: key.NewBinding(
				key.WithKeys("ctrl+c", "esc"),
				key.WithHelp("ctrl+c/esc", "exit"),
			),
			Backspace: key.NewBinding(
				key.WithKeys("backspace"),
				key.WithHelp("backspace", "delete character"),
			),
		},
		help: help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
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
		m.viewport.Height = msg.Height - m.textarea.Height() -
			lipgloss.Height(strings.Repeat("\n", 5))
		if len(m.outputs) > 0 {
			m.viewport.SetContent(
				lipgloss.NewStyle().
					Width(m.viewport.Width).
					Render(strings.Join(m.outputs, "\n")))
		}
		m.viewport.GotoBottom()

	case EvalCompleteMsg:
		return m.handleEvalComplete(msg, tiCmd, vpCmd)
	case ParserIoAvailable:
		return m.handleParserIoAvailable(msg, tiCmd, vpCmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Exit):
			if m.cancelFunc != nil {
				m.cancelFunc()
			}
			return m, tea.Quit
		case key.Matches(msg, m.keymap.ReverseSearch):
			// reverse search history
			return m.handleReverseHistorySearch(tiCmd, vpCmd)
		case key.Matches(msg, m.keymap.Submit):
			userCmd := m.pendingCommand.String()
			if userCmd == "" {
				return m, tea.Batch(tiCmd, vpCmd)
			}

			if m.input != "" {
				userCmd += " " + m.input + "\n"
			}
			if m.pendingCommand.Len() > 0 {
				m.pendingCommand.Reset()
			}
			m.input = ""
			return m.handleParserDispatch(userCmd, tiCmd, vpCmd)
		case key.Matches(msg, m.keymap.Enter):
			return m.handleLineOfCommandEntered(m.input, tiCmd, vpCmd)
		case key.Matches(msg, m.keymap.ScrollUpHistory):
			return m.handleScrollUpHistory(tiCmd, vpCmd)
		case key.Matches(msg, m.keymap.ScrollDownHistory):
			return m.handleScrollDownHistory(tiCmd, vpCmd)
		case key.Matches(msg, m.keymap.Backspace):
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

func (m Model) View() string {
	var historyStr strings.Builder
	//sep := lipgloss.NewStyle().BorderBottom(true).Width(m.textarea.Width())
	colsPerParen := 2 // spaces per open paren (tweak as needed)
	//indentCols := m.openParens * colsPerParen
	indentStyle := m.prevInputStyle
	historyStr.WriteString(m.headingStyle.Width(m.textarea.Width()).Render("------------- History --------") + "\n")
	var openParensCount int
	//m.renderOutputs(m.textarea.Width())
	indentCols := 0
	openParensCount = 0
	for _, cmd := range m.outputs {
		// apply a simple style in View as well
		if strings.Contains(cmd, "Executed:") {
			// keep the prefix (timestamp and "Executed:") and indent the command body
			idx := strings.Index(cmd, "Executed:")
			if idx >= 0 {
				prefix := cmd[:idx+len("Executed:")]
				body := strings.TrimSpace(cmd[idx+len("Executed:"):])
				if body != "" {
					openParensCount += strings.Count(body, "(") - strings.Count(body, ")")
					indentCols += openParensCount * colsPerParen

					historyStr.WriteString(prefix + "\n" + indentStyle.PaddingLeft(indentCols).Render(body) + "\n")
				} else {
					historyStr.WriteString(prefix + "\n")
				}
				continue
			}
		}

		if strings.Contains(cmd, "Output:") {
			historyStr.WriteString(m.outputStyle.Render(cmd) + "\n")
		} else {
			historyStr.WriteString(cmd + "\n")
		}
	}

	m.viewport.SetContent(historyStr.String())
	m.viewport.GotoBottom()

	inputStr := strings.Builder{}

	pending := strings.TrimSuffix(m.pendingCommand.String(), "\n")
	inputStr.WriteString(pending)
	if pending != "" {
		inputStr.WriteString("\n")
	}

	inputStr.WriteString(m.input)

	m.textarea.SetValue(inputStr.String())

	h := m.help.ShortHelpView([]key.Binding{
		m.keymap.Enter,
		//	m.keymap.ScrollUpHistory,
		//	m.keymap.ScrollDownHistory,
		//	m.keymap.ReverseSearch,
		m.keymap.Submit,
		m.keymap.Exit,
	})
	_ = historyStr.String()
	return lipgloss.JoinVertical(lipgloss.Center, m.viewport.View(), m.textarea.View()) + "\n" + h
}

func (m Model) handleEvalComplete(msg EvalCompleteMsg, tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	now := time.Now()
	if msg.err != nil {
		m.outputs = append(m.outputs,
			fmt.Sprintf("%s Error: %s", now.Format(time.Kitchen),
				msg.err.Error()))
	}

	m.viewport.SetContent(
		lipgloss.
			NewStyle().Width(m.viewport.Width).
			Render(strings.Join(m.outputs, "\n")))
	m.viewport.GotoBottom()
	m.input = ""
	m.cursor = -1
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) handleLineOfCommandEntered(msg string, tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	openParens := strings.Count(msg, "(")
	closeParens := strings.Count(msg, ")")
	m.openParens += openParens - closeParens

	m.history = append(m.history, msg)
	m.pendingCommand.WriteString(" ")
	m.pendingCommand.WriteString(msg)
	m.pendingCommand.WriteString("\n")

	m.input = ""
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) handleParserIoAvailable(msg ParserIoAvailable, tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	now := time.Now()
	if msg.content != "" {
		m.outputs = append(m.outputs,
			fmt.Sprintf("%s Output: %s",
				now.Format(time.Kitchen), msg.content))
		m.viewport.SetContent(lipgloss.NewStyle().
			Width(m.viewport.Width).
			Render(strings.Join(m.outputs, "\n")))
		m.viewport.GotoBottom()
	}
	m.input = ""
	m.cursor = -1
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) handleParserDispatch(msg string, tiCmd, vpCmd tea.Cmd) (tea.Model, tea.Cmd) {
	var (
		runParserCmd, updateParseResultCmd tea.Cmd
	)
	pr, pw := io.Pipe()
	runParserCmd = func() tea.Msg {
		m.runtime.Out = pw
		val, err := parser.EvalString(m.ctx, msg, m.runtime)
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
	m.outputs = append(m.outputs, fmt.Sprintf("%s Executed: %s", now.Format(time.Kitchen), msg))
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
