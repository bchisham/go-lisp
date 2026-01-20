package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

func NewVoidType() Interface {
	return Value{
		t: types.Void,
	}
}

func NewIdentifier(name string) Interface {
	return Value{
		t: types.Identifier,
		Primitive: Primitive{
			NameVal: name,
		},
	}
}

func NewString(s string) Interface {
	return Value{
		t: types.String,
		Primitive: Primitive{
			StringVal: s,
		},
	}
}

func NewInt(i int64) Interface {
	return Value{
		t: types.Int,
		Primitive: Primitive{
			IntVal:   i,
			FloatVal: float64(i),
		},
	}
}

func NewBool(b bool) Interface {
	return Value{
		t: types.Bool,
		Primitive: Primitive{
			BoolVal: b,
		},
	}
}

func NewChar(c rune) Interface {
	return Value{
		t: types.Char,
		Primitive: Primitive{
			CharVal: c,
		},
	}
}

func NewFloat(f float64) Interface {
	return Value{
		t: types.Float,
		Primitive: Primitive{
			FloatVal: f,
		},
	}
}
