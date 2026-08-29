package tui

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	Up                key.Binding
	Down              key.Binding
	Enter             key.Binding
	ReverseSearch     key.Binding
	ScrollUpHistory   key.Binding
	ScrollDownHistory key.Binding
	Submit            key.Binding
	Exit              key.Binding
	Backspace         key.Binding
}

type EvalCompleteMsg struct {
	result values.Type
	err    error
}

type ParserIoAvailable struct {
	content  string
	err      error
	canceled bool
}
