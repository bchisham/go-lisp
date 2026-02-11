package values

import (
	"fmt"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type String interface {
	Type
	fmt.Stringer
}

func NewString(value string) Type {
	return stringValue{value: value}
}

type stringValue struct {
	truthyValue
	value string
}

func (s stringValue) Equal(p Type) bool {
	if s.Type() != p.Type() {
		return false
	}
	other, ok := p.(stringValue)
	if !ok {
		return false
	}
	if s.value != other.value {
		return false
	}
	return true
}

func (s stringValue) Type() types.Type {
	return types.String
}

func (s stringValue) DisplayString() string {
	return fmt.Sprintf("%s", s.value)
}

func (s stringValue) WriteString() string {
	return fmt.Sprintf("%q", s.value)
}

func (s stringValue) String() string {
	return s.value
}
