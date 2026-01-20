package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

func NewQuotType() Interface {
	return Value{
		t: types.Quot,
	}
}

func NewRelationalOperator(literal string) Interface {
	return Value{
		t: types.RelationalOperator,
		Primitive: Primitive{
			Literal: literal,
		},
	}
}

func NewArithmeticOperator(literal string) Interface {
	return Value{
		t: types.ArithmeticOperator,
		Primitive: Primitive{
			Literal: literal,
		},
	}
}
