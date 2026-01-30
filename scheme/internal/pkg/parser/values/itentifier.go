package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

type Identifier interface {
	Interface
	GetName() string
}

func NewIdentifier(name string) Interface {
	return identifierValue{
		Literal: name,
	}
}

func (i identifierValue) GetName() string {
	return i.Literal
}

type identifierValue struct {
	truthyValue
	Literal string
}

func (i identifierValue) Equal(p Interface) bool {
	if i.Type() != p.Type() {
		return false
	}
	other, ok := p.(identifierValue)
	if !ok {
		return false
	}
	if i.Literal != other.Literal {
		return false
	}
	return true
}

func (i identifierValue) Type() types.Type {
	return types.Identifier
}

func (i identifierValue) DisplayString() string {
	return i.Literal
}

func (i identifierValue) WriteString() string {
	return i.Literal
}
