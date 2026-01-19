package parser

import (
	"slices"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

func defaultEnvironment() values.Environment {
	var builtins = values.NewEnvironment()
	builtins.Define("newline", values.NewString("\n"))
	//I/O
	builtins.Define(types.Format.String(), values.NewLambda(builtins, formatImpl))
	builtins.Define(types.Write.String(), values.NewLambda(builtins, writeImpl))
	builtins.Define(types.Display.String(), values.NewLambda(builtins, displayImpl))
	//relational operators
	builtins.Define("<", values.NewLambda(builtins, lessThanImpl))
	builtins.Define("<=", values.NewLambda(builtins, lessThanOrImpl))
	builtins.Define(">", values.NewLambda(builtins, greatThanImpl))
	builtins.Define(">=", values.NewLambda(builtins, greatThanOrImpl))
	builtins.Define("=", values.NewLambda(builtins, equalImpl))
	//boolean operators
	builtins.Define("not", values.NewLambda(builtins, notImpl))
	//arithmetic
	builtins.Define("+", values.NewLambda(builtins, sumImpl))
	builtins.Define("-", values.NewLambda(builtins, differenceImpl))
	builtins.Define("*", values.NewLambda(builtins, productImpl))

	return builtins
}

// evalSexpression evaluate a S-Expression i the given environment.
func evalSexpression(l []values.Value, env values.Environment) (values.Value, error) {
	head := list.Car(l)
	tail := list.Cdr(l)
	switch head.Type {
	case types.Quot:
		switch len(tail) {
		case 0:
			return values.NewVoidType(), ErrWrongNumberOfArguments
		case 1:
			return tail[0], nil
		default:
			return evalSexpression(tail, env)
		}
	case types.RelationalOperator, types.ArithmeticOperator:
		oper, ok := env.Lookup(head.Literal)
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		return oper.LambdaVal.Apply(tail)
	case types.Lambda:
		return head.LambdaVal.Apply(tail)
	case types.Identifier:
		resvVal, ok := env.Lookup(head.NameVal)
		if !ok {
			return values.NewVoidType(), ErrUndefinedIdent
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
		if len(tail) == 0 {
			return head, nil
		}
		return values.NewList(l...), nil
	}
}

func quotImpl(args []values.Value, env values.Environment) (values.Value, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	return values.NewList(args...), nil
}

func relationalCompareAllowed(t types.Type) bool {
	return slices.Contains(list.New(types.Int, types.Float), t)
}

func lessThanImpl(args []values.Value, env values.Environment) (values.Value, error) {
	//https://try.scheme.org/ returns #t when there are no operands
	if len(args) == 0 {
		return values.NewBool(true), nil
	}

	if len(args) == 1 {
		return values.NewBool(relationalCompareAllowed(args[0].Type)), nil
	}
	head := list.Car(args)
	tail := list.Cdr(args)
	second := list.Car(tail)

	if !relationalCompareAllowed(head.Type) || !relationalCompareAllowed(second.Type) {
		return values.NewBool(false), nil
	}
	tailIsInvariant, err := lessThanImpl(tail, env)
	if err != nil {
		return values.NewBool(false), err
	}

	return values.NewBool(head.FloatVal < second.FloatVal && tailIsInvariant.BoolVal), nil

}

func lessThanOrImpl(args []values.Value, env values.Environment) (values.Value, error) {
	if len(args) == 0 {
		return values.NewBool(true), nil
	}
	if len(args) == 1 {
		return values.NewBool(relationalCompareAllowed(args[0].Type)), nil
	}
	head := list.Car(args)
	tail := list.Cdr(args)
	second := list.Car(tail)
	if !relationalCompareAllowed(head.Type) || !relationalCompareAllowed(second.Type) {
		return values.NewBool(false), nil
	}
	tailIsInvariant, err := lessThanOrImpl(tail, env)
	if err != nil {
		return values.NewBool(false), err
	}
	return values.NewBool(head.FloatVal <= second.FloatVal && tailIsInvariant.BoolVal), nil
}

func greatThanImpl(args []values.Value, env values.Environment) (values.Value, error) {
	if len(args) == 0 {
		return values.NewBool(true), nil
	}
	if len(args) == 1 {
		return values.NewBool(relationalCompareAllowed(args[0].Type)), nil
	}
	head := list.Car(args)
	tail := list.Cdr(args)
	second := list.Car(tail)
	if !relationalCompareAllowed(head.Type) || !relationalCompareAllowed(second.Type) {
		return values.NewBool(false), nil
	}
	tailIsInvariant, err := greatThanImpl(tail, env)
	if err != nil {
		return values.NewBool(false), err
	}
	return values.NewBool(head.FloatVal > second.FloatVal && tailIsInvariant.BoolVal), nil
}

func greatThanOrImpl(args []values.Value, env values.Environment) (values.Value, error) {
	if len(args) == 0 {
		return values.NewBool(true), nil
	}
	if len(args) == 1 {
		return values.NewBool(relationalCompareAllowed(args[0].Type)), nil
	}
	head := list.Car(args)
	tail := list.Cdr(args)
	second := list.Car(tail)
	if !relationalCompareAllowed(head.Type) || !relationalCompareAllowed(second.Type) {
		return values.NewBool(false), nil
	}
	tailIsInvariant, err := greatThanOrImpl(tail, env)
	if err != nil {
		return values.NewBool(false), err
	}
	return values.NewBool(head.FloatVal >= second.FloatVal && tailIsInvariant.BoolVal), nil
}

func equalImpl(args []values.Value, env values.Environment) (values.Value, error) {
	if len(args) == 0 {
		return values.NewBool(true), nil
	}
	if len(args) == 1 {
		return values.NewBool(relationalCompareAllowed(args[0].Type)), nil
	}
	head := list.Car(args)
	tail := list.Cdr(args)
	second := list.Car(tail)
	if !relationalCompareAllowed(head.Type) || !relationalCompareAllowed(second.Type) {
		return values.NewBool(false), nil
	}
	tailIsInvariant, err := equalImpl(tail, env)
	if err != nil {
		return values.NewBool(false), err
	}
	return values.NewBool(head.FloatVal == second.FloatVal && tailIsInvariant.BoolVal), nil
}

func notImpl(args []values.Value, env values.Environment) (values.Value, error) {
	if len(args) != 1 {
		return values.NewBool(false), ErrWrongNumberOfArguments
	}
	head := list.Car(args)
	switch head.Type {
	case types.Bool:
		return values.NewBool(!head.BoolVal), nil
	default:
		return values.NewBool(false), nil
	}
}

var arithmeticAllowedTypes = list.New(types.Int, types.Float)

func arithmeticAllowed(t types.Type) bool {
	return slices.Contains(arithmeticAllowedTypes, t)
}

func sumImpl(args []values.Value, env values.Environment) (values.Value, error) {
	if len(args) == 0 {
		return values.NewInt(0), nil
	}
	if len(args) == 1 {
		if !arithmeticAllowed(args[0].Type) {
			return values.NewVoidType(), ErrNumberExpected
		}
		return args[0], nil
	}

	sum := values.NewInt(0)
	for i := 0; i < len(args); i++ {
		if arithmeticAllowed(args[i].Type) {
			sum.FloatVal += args[i].FloatVal
			sum.IntVal += args[i].IntVal
		} else {
			return values.NewVoidType(), ErrNumberExpected
		}
	}
	return sum, nil
}

func differenceImpl(args []values.Value, env values.Environment) (values.Value, error) {
	if len(args) == 0 {
		return values.NewInt(0), nil
	}
	if len(args) == 1 {
		if !arithmeticAllowed(args[0].Type) {
			return values.NewVoidType(), ErrNumberExpected
		}
		return args[0], nil
	}
	diff := values.NewInt(0)
	for i := 0; i < len(args); i++ {
		if arithmeticAllowed(args[i].Type) {
			diff.FloatVal -= args[i].FloatVal
			diff.IntVal -= args[i].IntVal
		} else {
			return values.NewVoidType(), ErrNumberExpected
		}
	}
	return diff, nil
}

func productImpl(args []values.Value, env values.Environment) (values.Value, error) {
	if len(args) == 0 {
		return values.NewInt(0), nil
	}
	if len(args) == 1 {
		if !arithmeticAllowed(args[0].Type) {
			return values.NewVoidType(), ErrNumberExpected
		}
		return args[0], nil
	}
	product := values.NewInt(1)
	for i := 0; i < len(args); i++ {
		if arithmeticAllowed(args[i].Type) {
			product.FloatVal *= args[i].FloatVal
			product.IntVal *= args[i].IntVal
		} else {
			return values.NewVoidType(), ErrNumberExpected
		}
	}
	return product, nil
}
