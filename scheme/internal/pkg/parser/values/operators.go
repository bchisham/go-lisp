package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

func NewQuotType() Value {
	return Value{
		Type: types.Quot,
	}
}

func NewRelationalOperator(literal string) Value {
	return Value{
		Type: types.RelationalOperator,
		Primitive: Primitive{
			Literal: literal,
		},
	}
}

func NewArithmeticOperator(literal string) Value {
	return Value{
		Type: types.ArithmeticOperator,
		Primitive: Primitive{
			Literal: literal,
		},
	}
}
