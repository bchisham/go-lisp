package values

import (
	"fmt"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type Char interface {
	Interface
	Rune() rune
}

type char struct {
	truthyValue
	rune rune
}

func (c char) Equal(p Interface) bool {
	//TODO implement me
	panic("implement me")
}

func (c char) Type() types.Type {
	return types.Char
}

func (c char) DisplayString() string {
	return fmt.Sprintf("%q", c.rune)
}

func (c char) WriteString() string {
	return writeChar(c.rune)
}

func NewChar(r rune) Interface {
	return char{rune: r}
}

func writeChar(r rune) string {
	switch r {
	case '\n':
		return "#\\newline"
	case ' ':
		return "#\\space"
	case '\t':
		return "#\\tab"
	// add others if you like
	default:
		// If you want #\a for printable runes:
		return "#\\" + string(r)
	}
}
