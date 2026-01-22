package values

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	tokentypes "github.com/bchisham/go-lisp/scheme/internal/pkg/lexer/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type NilValue struct{}

func NewNil() Interface {
	return NilValue{}
}

func (n NilValue) String() string {
	return "()"
}

func (n NilValue) Equal(p Interface) bool {
	_, ok := p.(NilValue)
	return ok
}

func (n NilValue) Type() types.Type {
	return types.Nil
}

func (n NilValue) GetToken() lexer.Token {
	return lexer.Token{
		Type:    tokentypes.LiteralNilList,
		Literal: "()",
	}
}

func (n NilValue) AsPrimitive() (Primitive, error) {
	return Primitive{}, ErrNotAPrimitive
}

func (n NilValue) IsTruthy() bool {
	return true
}

func (n NilValue) DisplayString() string {
	return "()"
}

func (n NilValue) WriteString() string {
	return "()"
}
