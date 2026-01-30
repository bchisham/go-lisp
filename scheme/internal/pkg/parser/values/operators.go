package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

type Operator interface {
	Interface
	IsArithmetic() bool
	IsRelational() bool
	IsBoolean() bool
	GetName() string
}

type operatorValue struct {
	truthyValue
	t       types.Type
	Literal string
}

func (o operatorValue) IsArithmetic() bool {
	return o.Type() == types.ArithmeticOperator
}

func (o operatorValue) IsRelational() bool {
	return o.Type() == types.RelationalOperator
}

func (o operatorValue) IsBoolean() bool {
	return o.Type() == types.BooleanOperator
}

func (o operatorValue) GetName() string {
	return o.Literal
}

func (o operatorValue) Equal(p Interface) bool {
	if o.Type() != p.Type() {
		return false
	}
	other, ok := p.(operatorValue)
	if !ok {
		return false
	}
	if o.Literal != other.Literal {
		return false
	}
	return true
}

func (o operatorValue) Type() types.Type {
	return o.t
}

func (o operatorValue) DisplayString() string {
	return o.Literal
}

func (o operatorValue) WriteString() string {
	return o.Literal
}

func NewRelationalOperator(literal string) Interface {
	return operatorValue{
		t:       types.RelationalOperator,
		Literal: literal,
	}
}

func NewArithmeticOperator(literal string) Interface {
	return operatorValue{
		t:       types.ArithmeticOperator,
		Literal: literal,
	}
}
