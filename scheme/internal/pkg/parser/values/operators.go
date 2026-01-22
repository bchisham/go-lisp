package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

type operatorValue struct {
	t       types.Type
	Literal string
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

func (o operatorValue) AsPrimitive() (Primitive, error) {
	return Primitive{
		Literal: o.Literal,
	}, nil
}

func (o operatorValue) IsTruthy() bool {
	return true
}

func (o operatorValue) DisplayString() string {
	return o.Literal
}

func (o operatorValue) WriteString() string {
	return o.Literal
}

func NewQuotType() Interface {
	return operatorValue{
		t: types.Quot,
	}
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
