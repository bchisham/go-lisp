package parser

import (
	"errors"
	"slices"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

// evalSexpression evaluate a S-Expression i the given environment.
func evalSexpression(l values.Interface, rt *Runtime) (values.Interface, error) {

	switch l.Type() {
	case types.Nil:
		return values.NewNil(), nil
	case types.ArithmeticOperator, types.RelationalOperator:
		p, err := l.AsPrimitive()
		if err != nil {
			return values.NewVoidType(), err
		}
		oper, ok := rt.Env.Lookup(p.Literal)
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		lambda, ok := oper.(LambdaExpr)
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		return lambda.Apply(values.NewNil())
	case types.Lambda:
		lambda, ok := l.(LambdaExpr)
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		return lambda.Apply(values.NewNil())
	case types.Quot:
		return values.NewNil(), nil
	case types.Identifier:
		p, err := l.AsPrimitive()
		if err != nil {
			return values.NewVoidType(), err
		}
		resolvVal, ok := rt.Env.Lookup(p.NameVal)
		if !ok {
			return values.NewVoidType(), ErrUndefinedIdent
		}
		return resolvVal, nil
	}

	head := values.Car(l)
	tail := values.Cdr(l)
	switch head.Type() {
	case types.Pair:
		evaluatedHead, err := evalSexpression(head, rt)
		if err != nil && !errors.Is(err, ErrOperatorIsNotAProcedure) {
			return evaluatedHead, err
		}
		switch evaluatedHead.Type() {
		case types.Identifier:
			p, err := evaluatedHead.AsPrimitive()
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
						return values.NewNil(), ErrOperatorIsNotAProcedure
					}
					return lambda.Apply(tail)
				}
				return resvVal, nil

			}
		case types.RelationalOperator, types.ArithmeticOperator:
			p, err := evaluatedHead.AsPrimitive()
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
			lambda, ok := evaluatedHead.(LambdaExpr)
			if !ok {
				return values.NewVoidType(), ErrOperatorIsNotAProcedure
			}
			return lambda.Apply(tail)
		case types.Quot:
			return tail, nil
		default:
			return evaluatedHead, ErrOperatorIsNotAProcedure
		}
	case types.Quot:
		switch tail.Type() {
		case types.Nil:
			return values.NewNil(), ErrWrongNumberOfArguments
		default:
			return tail, nil
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
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		return lambda.Apply(tail)
	case types.Lambda:
		lambda, ok := head.(LambdaExpr)
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		return lambda, nil

	case types.Identifier:
		p, err := head.AsPrimitive()
		if err != nil {
			return values.NewVoidType(), err
		}
		resolvVal, ok := rt.Env.Lookup(p.NameVal)
		if !ok {
			return values.NewVoidType(), ErrUndefinedIdent
		}
		if resolvVal.Type() == types.Lambda {
			lambda, ok := resolvVal.(LambdaExpr)
			if !ok {
				return values.NewNil(), ErrOperatorIsNotAProcedure
			}
			return lambda.Apply(tail)
		}
		return resolvVal, nil
	default:
		if tail.Type() == types.Nil {
			return head, nil
		}

	}
	return l, ErrOperatorIsNotAProcedure
}

func quotImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Pair {
		return values.Car(args), nil
	}
	return args, nil
}

func relationalCompareAllowed(t types.Type) bool {
	return slices.Contains(list.New(types.Int, types.Float), t)
}

func lessThanImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	//https://try.scheme.org/ returns #t when there are no operands
	if args.Type() == types.Nil {
		return values.NewBool(true), nil
	}
	head := values.Car(args)
	tail := values.Cdr(args)
	if tail.Type() == types.Nil {
		return values.NewBool(relationalCompareAllowed(head.Type())), nil
	}

	second := values.Car(tail)

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

func lessThanOrImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewBool(true), nil
	}
	head := values.Car(args)
	tail := values.Cdr(args)
	if tail.Type() == types.Nil {
		return values.NewBool(relationalCompareAllowed(head.Type())), nil
	}

	second := values.Car(tail)
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

func greatThanImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewBool(true), nil
	}
	head := values.Car(args)
	tail := values.Cdr(args)
	if tail.Type() == types.Nil {
		return values.NewBool(relationalCompareAllowed(head.Type())), nil
	}

	second := values.Car(tail)
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

func greatThanOrImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewBool(true), nil
	}
	head := values.Car(args)
	tail := values.Cdr(args)
	if tail.Type() == types.Nil {
		return values.NewBool(relationalCompareAllowed(head.Type())), nil
	}

	second := values.Car(tail)
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

func equalImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewBool(true), nil
	}
	head := values.Car(args)
	tail := values.Cdr(args)
	if tail.Type() == types.Nil {
		return values.NewBool(relationalCompareAllowed(head.Type())), nil
	}

	second := values.Car(tail)
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

func notImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewBool(false), ErrWrongNumberOfArguments
	}
	return values.NewBool(!values.Car(args).IsTruthy()), nil
}

var arithmeticAllowedTypes = list.New(types.Int, types.Float)

func arithmeticAllowed(t types.Type) bool {
	return slices.Contains(arithmeticAllowedTypes, t)
}

func sumImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewInt(0), nil
	}
	head := values.Car(args)
	tail := values.Cdr(args)
	sum := values.Primitive{}
	if tail.Type() == types.Nil {
		if !arithmeticAllowed(head.Type()) {
			return values.NewNil(), ErrNumberExpected
		}
		return head, nil
	}
	ti, err := sumImpl(tail, rt)
	if err != nil {
		return values.NewNil(), err
	}
	tailSum, err := ti.AsPrimitive()
	if err != nil {
		return values.NewNil(), err
	}
	lhs, err := head.AsPrimitive()
	if err != nil {
		return values.NewNil(), err
	}
	if !arithmeticAllowed(head.Type()) {
		return values.NewNil(), ErrNumberExpected
	}
	sum.FloatVal = lhs.FloatVal + tailSum.FloatVal
	sum.IntVal = lhs.IntVal + tailSum.IntVal

	return values.FromPrimitive(types.Int, sum), nil
}

func differenceImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewInt(0), nil
	}
	head := values.Car(args)
	tail := values.Cdr(args)
	if tail.Type() == types.Nil {
		if !arithmeticAllowed(head.Type()) {
			return values.NewNil(), ErrNumberExpected
		}
		return head, nil
	}
	diff := values.Primitive{}
	lhs, err := head.AsPrimitive()
	if err != nil {
		return values.NewNil(), err
	}
	second := values.Car(tail)
	if !arithmeticAllowed(head.Type()) || !arithmeticAllowed(second.Type()) {
		return values.NewNil(), ErrNumberExpected
	}
	rhs, err := second.AsPrimitive()
	if err != nil {
		return values.NewNil(), err
	}
	diff.FloatVal = lhs.FloatVal - rhs.FloatVal
	diff.IntVal = lhs.IntVal - rhs.IntVal

	ti, err := differenceImpl(tail, rt)
	if err != nil {
		return values.NewNil(), err
	}
	tailDiff, err := ti.AsPrimitive()
	if err != nil {
		return values.NewNil(), err
	}
	diff.FloatVal -= tailDiff.FloatVal
	diff.IntVal -= tailDiff.IntVal
	return values.FromPrimitive(types.Int, diff), nil
}

func productImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewInt(0), nil
	}
	head := values.Car(args)
	tail := values.Cdr(args)
	if tail.Type() == types.Nil {
		if !arithmeticAllowed(head.Type()) {
			return values.NewNil(), ErrNumberExpected
		}
		return head, nil
	}
	product := values.Primitive{
		IntVal:   1,
		FloatVal: 1,
	}
	lhs, err := head.AsPrimitive()
	if err != nil {
		return values.NewNil(), err
	}
	second := values.Car(tail)
	if !arithmeticAllowed(head.Type()) || !arithmeticAllowed(second.Type()) {
		return values.NewNil(), ErrNumberExpected
	}
	rhs, err := second.AsPrimitive()
	if err != nil {
		return values.NewNil(), err
	}
	product.FloatVal = lhs.FloatVal * rhs.FloatVal
	product.IntVal = lhs.IntVal * rhs.IntVal

	ti, err := productImpl(tail, rt)
	if err != nil {
		return values.NewNil(), err
	}
	tailProduct, err := ti.AsPrimitive()
	if err != nil {
		return values.NewNil(), err
	}
	product.FloatVal *= tailProduct.FloatVal
	product.IntVal *= tailProduct.IntVal
	return values.FromPrimitive(types.Int, product), nil
}
