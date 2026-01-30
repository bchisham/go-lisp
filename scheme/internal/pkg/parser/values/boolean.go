package values

import (
	"fmt"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/boolean"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type Boolean interface {
	Interface
	GetLiteral() bool
}

func NewBool(b bool) Interface {
	return booleanValue{
		Literal: b,
	}
}

type booleanValue struct {
	Literal bool
}

func (b booleanValue) Equal(p Interface) bool {
	if b.Type() != p.Type() {
		return false
	}
	other, ok := p.(booleanValue)
	if !ok {
		return false
	}
	if b.Literal != other.Literal {
		return false
	}
	return true
}

func (b booleanValue) Type() types.Type {
	return types.Bool
}

func (b booleanValue) IsTruthy() bool {
	return b.Literal
}

func (b booleanValue) DisplayString() string {
	return fmt.Sprintf("%v", boolean.Trinary(b.Literal, "#t", "#f"))
}

func (b booleanValue) WriteString() string {
	return fmt.Sprintf("%v", boolean.Trinary(b.Literal, "#t", "#f"))
}

func (b booleanValue) GetLiteral() bool {
	return b.Literal
}
