package types

import (
	"fmt"
	"strconv"

	"github.com/bchisham/go-lisp/scheme/internal/contracts"
	"github.com/charmbracelet/lipgloss"
)

type stringCanDisplay struct {
	displayConfig
	val string
}
type intCanDisplay struct {
	displayConfig
	val int
}

type displayConfig struct {
	fmt   string
	base  int
	style lipgloss.Style
}

type anyCanDisplay[T any] struct {
	displayConfig
	val T
}

type DisplayConfigOption func(*displayConfig)

func WithDisplayFormat(fmt string) DisplayConfigOption {
	return func(c *displayConfig) {
		c.fmt = fmt
	}
}

func WithDisplayBase(base int) DisplayConfigOption {
	return func(c *displayConfig) {
		c.base = base
	}
}

func WithDisplayStyle(style lipgloss.Style) DisplayConfigOption {
	return func(c *displayConfig) {
		c.style = style
	}
}

func NewStringCanDisplay(s string, opts ...DisplayConfigOption) contracts.CanDisplay {
	cfg := &displayConfig{
		fmt:  "%v",
		base: 10,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return stringCanDisplay{
		val:           s,
		displayConfig: *cfg,
	}
}

func NewAnyCanDisplay[T any](val T, opts ...DisplayConfigOption) contracts.CanDisplay {
	cfg := &displayConfig{
		fmt:  "%v",
		base: 10,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return anyCanDisplay[T]{
		val:           val,
		displayConfig: *cfg,
	}
}

func NewAnyCanDisplayFunc[T any](opts ...DisplayConfigOption) func(t T) contracts.CanDisplay {
	cfg := &displayConfig{
		fmt:  "%v",
		base: 10,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(t T) contracts.CanDisplay {
		return anyCanDisplay[T]{
			val:           t,
			displayConfig: *cfg,
		}
	}
}

func (s stringCanDisplay) DisplayString() string {
	return s.style.Render(fmt.Sprintf(s.fmt, s.val))
}

func NewIntCanDisplay(i int, opts ...DisplayConfigOption) contracts.CanDisplay {
	cfg := &displayConfig{
		fmt:  "%v",
		base: 10,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return intCanDisplay{
		val:           i,
		displayConfig: *cfg,
	}
}

func (i intCanDisplay) DisplayString() string {
	return i.style.Render(strconv.Itoa(i.val))
}

func (a anyCanDisplay[T]) DisplayString() string {
	return a.style.Render(fmt.Sprintf(a.fmt, a.val))
}
