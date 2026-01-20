package parser

import (
	"fmt"
	"slices"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

// evalSexpression evaluate a S-Expression i the given environment.
func evalSexpression(l []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(l) == 0 {
		return values.NewVoidType(), fmt.Errorf("no expression found")
	}
	head := list.Car(l)
	tail := list.Cdr(l)
	switch head.Type() {
	case types.Quot:
		switch len(tail) {
		case 0:
			return values.NewVoidType(), ErrWrongNumberOfArguments
		case 1:
			return tail[0], nil
		default:
			return evalSexpression(tail, rt)
		}
	case types.RelationalOperator, types.ArithmeticOperator:
		p, err := head.AsPrimitive()
		if err != nil {
			return values.NewVoidType(), err
		}
		oper, ok := rt.Env.Lookup(p.Literal)
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		lambda, ok := oper.(LambdaExpr)

		return lambda.Apply(tail)
	case types.Lambda:
		lambda, ok := head.(LambdaExpr)
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		return lambda.Apply(tail)
	case types.Identifier:
		p, err := head.AsPrimitive()
		if err != nil {
			return values.NewVoidType(), err
		}
		resvVal, ok := rt.Env.Lookup(p.NameVal)
		if !ok {
			return values.NewVoidType(), ErrUndefinedIdent
		}
		_, userDefined := types.FromString(p.NameVal)
		switch resvVal.Type() {
		case types.Lambda:
			if userDefined == nil {
				lambda, ok := resvVal.(LambdaExpr)
				if !ok {
					return values.NewVoidType(), ErrOperatorIsNotAProcedure
				}
				return lambda.Apply(tail)
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

func quotImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	return values.NewList(args...), nil
}

func relationalCompareAllowed(t types.Type) bool {
	return slices.Contains(list.New(types.Int, types.Float), t)
}

func lessThanImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	//https://try.scheme.org/ returns #t when there are no operands
	if len(args) == 0 {
		return values.NewBool(true), nil
	}

	if len(args) == 1 {
		return values.NewBool(relationalCompareAllowed(args[0].Type())), nil
	}
	head := list.Car(args)
	tail := list.Cdr(args)
	second := list.Car(tail)

	if !relationalCompareAllowed(head.Type()) ||
		!relationalCompareAllowed(second.Type()) {
		return values.NewBool(false), nil
	}
	ti, err := lessThanImpl(tail, rt)
	if err != nil {
		return values.NewBool(false), err
	}
	tailIsInvariant, err := ti.AsPrimitive()
	if err != nil {
		return values.NewBool(false), nil
	}

	lhs, err := head.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}
	rhs, err := second.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}

	return values.NewBool(lhs.FloatVal < rhs.FloatVal && tailIsInvariant.BoolVal), nil

}

func lessThanOrImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(args) == 0 {
		return values.NewBool(true), nil
	}
	if len(args) == 1 {
		return values.NewBool(relationalCompareAllowed(args[0].Type())), nil
	}
	head := list.Car(args)
	tail := list.Cdr(args)
	second := list.Car(tail)
	if !relationalCompareAllowed(head.Type()) || !relationalCompareAllowed(second.Type()) {
		return values.NewBool(false), nil
	}
	ti, err := lessThanOrImpl(tail, rt)
	if err != nil {
		return values.NewBool(false), err
	}
	tailIsInvariant, err := ti.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}
	lhs, err := head.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}
	rhs, err := second.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}

	return values.NewBool(lhs.FloatVal <= rhs.FloatVal && tailIsInvariant.BoolVal), nil
}

func greatThanImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(args) == 0 {
		return values.NewBool(true), nil
	}
	if len(args) == 1 {
		return values.NewBool(relationalCompareAllowed(args[0].Type())), nil
	}
	head := list.Car(args)
	tail := list.Cdr(args)
	second := list.Car(tail)
	if !relationalCompareAllowed(head.Type()) ||
		!relationalCompareAllowed(second.Type()) {
		return values.NewBool(false), nil
	}
	ti, err := greatThanImpl(tail, rt)
	if err != nil {
		return values.NewBool(false), err
	}
	tailIsInvariant, err := ti.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}
	lhs, err := head.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}
	rhs, err := second.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}

	return values.NewBool(lhs.FloatVal > rhs.FloatVal && tailIsInvariant.BoolVal), nil
}

func greatThanOrImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(args) == 0 {
		return values.NewBool(true), nil
	}
	if len(args) == 1 {
		return values.NewBool(relationalCompareAllowed(args[0].Type())), nil
	}
	head := list.Car(args)
	tail := list.Cdr(args)
	second := list.Car(tail)
	if !relationalCompareAllowed(head.Type()) || !relationalCompareAllowed(second.Type()) {
		return values.NewBool(false), nil
	}
	ti, err := greatThanOrImpl(tail, rt)
	if err != nil {
		return values.NewBool(false), err
	}
	tailIsInvariant, err := ti.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}
	lhs, err := head.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}
	rhs, err := second.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}

	return values.NewBool(lhs.FloatVal >= rhs.FloatVal && tailIsInvariant.BoolVal), nil
}

func equalImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(args) == 0 {
		return values.NewBool(true), nil
	}
	if len(args) == 1 {
		return values.NewBool(relationalCompareAllowed(args[0].Type())), nil
	}
	head := list.Car(args)
	tail := list.Cdr(args)
	second := list.Car(tail)
	if !relationalCompareAllowed(head.Type()) || !relationalCompareAllowed(second.Type()) {
		return values.NewBool(false), nil
	}
	ti, err := equalImpl(tail, rt)
	if err != nil {
		return values.NewBool(false), err
	}

	tailIsInvariant, err := ti.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}
	lhs, err := head.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}

	rhs, err := second.AsPrimitive()
	if err != nil {
		return values.NewBool(false), err
	}
	return values.NewBool(lhs.FloatVal == rhs.FloatVal && tailIsInvariant.BoolVal), nil
}

func notImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(args) != 1 {
		return values.NewBool(false), ErrWrongNumberOfArguments
	}
	head := list.Car(args)
	switch head.Type() {
	case types.Bool:
		h, err := head.AsPrimitive()
		if err != nil {
			return values.NewBool(false), err
		}
		return values.NewBool(!h.BoolVal), nil
	default:
		return values.NewBool(false), nil
	}
}

var arithmeticAllowedTypes = list.New(types.Int, types.Float)

func arithmeticAllowed(t types.Type) bool {
	return slices.Contains(arithmeticAllowedTypes, t)
}

func sumImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(args) == 0 {
		return values.NewInt(0), nil
	}
	if len(args) == 1 {
		if !arithmeticAllowed(args[0].Type()) {
			return values.NewVoidType(), ErrNumberExpected
		}
		return args[0], nil
	}

	sum := values.Primitive{}

	for i := 0; i < len(args); i++ {
		if arithmeticAllowed(args[i].Type()) {
			rhs, err := args[i].AsPrimitive()
			if err != nil {
				return values.NewInt(0), err
			}
			sum.FloatVal += rhs.FloatVal
			sum.IntVal += rhs.IntVal
		} else {
			return values.NewVoidType(), ErrNumberExpected
		}
	}
	return values.FromPrimitive(types.Int, sum), nil
}

func differenceImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(args) == 0 {
		return values.NewInt(0), nil
	}
	if len(args) == 1 {
		if !arithmeticAllowed(args[0].Type()) {
			return values.NewVoidType(), ErrNumberExpected
		}
		return args[0], nil
	}
	diff := values.Primitive{}
	for i := 0; i < len(args); i++ {
		if arithmeticAllowed(args[i].Type()) {
			rhs, err := args[i].AsPrimitive()
			if err != nil {
				return values.NewInt(0), err
			}
			diff.FloatVal -= rhs.FloatVal
			diff.IntVal -= rhs.IntVal
		} else {
			return values.NewVoidType(), ErrNumberExpected
		}
	}
	return values.FromPrimitive(types.Int, diff), nil
}

func productImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(args) == 0 {
		return values.NewInt(0), nil
	}
	if len(args) == 1 {
		if !arithmeticAllowed(args[0].Type()) {
			return values.NewVoidType(), ErrNumberExpected
		}
		return args[0], nil
	}
	product := values.Primitive{
		IntVal:   1,
		FloatVal: 1,
	}
	for i := 0; i < len(args); i++ {
		if arithmeticAllowed(args[i].Type()) {
			rhs, err := args[i].AsPrimitive()
			if err != nil {
				return values.NewInt(0), err
			}
			product.FloatVal *= rhs.FloatVal
			product.IntVal *= rhs.IntVal
		} else {
			return values.NewVoidType(), ErrNumberExpected
		}
	}
	return values.FromPrimitive(types.Int, product), nil
}
