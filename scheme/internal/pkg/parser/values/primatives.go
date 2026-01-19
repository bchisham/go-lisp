package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

func NewVoidType() Value {
	return Value{
		Type: types.Void,
	}
}

func NewIdentifier(name string) Value {
	return Value{
		Type: types.Identifier,
		Primitive: Primitive{
			NameVal: name,
		},
	}
}

func NewString(s string) Value {
	return Value{
		Type: types.String,
		Primitive: Primitive{
			StringVal: s,
		},
	}
}

func NewInt(i int64) Value {
	return Value{
		Type: types.Int,
		Primitive: Primitive{
			IntVal:   i,
			FloatVal: float64(i),
		},
	}
}

func NewBool(b bool) Value {
	return Value{
		Type: types.Bool,
		Primitive: Primitive{
			BoolVal: b,
		},
	}
}

func NewChar(c rune) Value {
	return Value{
		Type: types.Char,
		Primitive: Primitive{
			CharVal: string(c),
		},
	}
}

func NewFloat(f float64) Value {
	return Value{
		Type: types.Float,
		Primitive: Primitive{
			FloatVal: f,
		},
	}
}
