package parser

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

func defaultEnvironment() types.Environment {
	var builtins = types.NewEnvironment()
	builtins.Define("newline", newString("\n"))
	builtins.Define(types.Format.String(), newLambda(builtins, formatImpl))
	builtins.Define(types.Write.String(), newLambda(builtins, writeImpl))
	builtins.Define(types.Display.String(), newLambda(builtins, displayImpl))
	builtins.Define(types.Quot.String(), newLambda(builtins, quotImpl))

	return builtins
}

func evalSexpression(l []types.Value, env types.Environment) (types.Value, error) {
	head := list.Car(l)
	tail := list.Cdr(l)
	switch head.Type {
	case types.Lambda:
		return head.LambdaVal.Apply(tail)
	case types.Identifier:
		resvVal, ok := env.Lookup(head.NameVal)
		if !ok {
			return newVoidType(), ErrUndefinedIdent
		}
		_, userDefined := types.FromString(head.NameVal)
		switch resvVal.Type {
		case types.Lambda:
			if userDefined == nil {
				return resvVal.LambdaVal.Apply(tail)
			}
			return resvVal, nil

		default:
			return resvVal, nil
		}

	default:
		return newList(l...), nil
	}
}

func quotImpl(args []types.Value, env types.Environment) (types.Value, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	return newList(args...), nil
}
