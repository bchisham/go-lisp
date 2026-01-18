package parser

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

func newVoidType() types.Value {
	return types.Value{
		Type: types.Void,
	}
}

func newString(s string) types.Value {
	return types.Value{
		Type: types.String,
		Primitive: types.Primitive{
			StringVal: s,
		},
	}
}

func newFloat(f float64) types.Value {
	return types.Value{
		Type: types.Float,
		Primitive: types.Primitive{
			FloatVal: f,
		},
	}
}

func newInt(i int64) types.Value {
	return types.Value{
		Type: types.Int,
		Primitive: types.Primitive{
			IntVal: i,
		},
	}
}

func newBool(b bool) types.Value {
	return types.Value{
		Type: types.Bool,
		Primitive: types.Primitive{
			BoolVal: b,
		},
	}
}
func newChar(c rune) types.Value {
	return types.Value{
		Type: types.Char,
		Primitive: types.Primitive{
			CharVal: string(c),
		},
	}
}

func newLambda(env types.Environment, expression types.Expression) types.Value {
	return types.Value{
		Type: types.Lambda,
		LambdaVal: types.LambdaExpr{
			Env:  env,
			Body: expression,
		},
	}
}

func newExpression(env types.Environment, body []types.Value) types.Expression {
	return func(args []types.Value, env types.Environment) (types.Value, error) {
		//TODO evaluate body in environment
		return newVoidType(), nil
	}
}

func newIdentifier(name string) types.Value {
	return types.Value{
		Type: types.Identifier,
		Primitive: types.Primitive{
			NameVal: name,
		},
	}
}

func newList(v ...types.Value) types.Value {
	if len(v) == 0 {
		return types.Value{
			Type:    types.List,
			ListVal: []types.Value{},
		}
	}
	if len(v) == 1 && v[0].Type == types.List {
		return v[0]
	}
	return types.Value{
		Type:    types.List,
		ListVal: v,
	}
}
