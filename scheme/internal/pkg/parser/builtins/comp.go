package builtins

import (
	"slices"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

func relationalCompareAllowed(t types.Type) bool {
	return slices.Contains(list.New(types.Int, types.Float), t)
}

func LessThanImpl(args values.Interface, rt *Runtime, cb Expression) (_ values.Interface, err error) {
	//https://try.scheme.org/ returns #t when there are no operands
	if args.Type() == types.Nil {
		return values.NewBool(true), nil
	}
	if relationalCompareAllowed(args.Type()) {
		return values.NewBool(true), nil
	}
	if args.Type() != types.Pair {
		return values.NewBool(false), ErrWrongNumberOfArguments
	}
	l := args.(values.Pair)
	head := l.Car()
	var (
		lhs values.Numeric
		rhs values.Numeric
	)
	lhs, err = evalToNumber(head, rt, cb)
	if err != nil {
		return values.NewBool(false), err
	}
	var invariant = lhs != nil

	tail := l.Cdr()

	for tail.Type() != types.Nil && invariant {
		current := values.Car(tail)
		rhs, err = evalToNumber(current, rt, cb)
		if err != nil {
			return values.NewBool(false), err
		}
		invariant = rhs != nil && lhs.LessThan(rhs)
		lhs = rhs
		tail = values.Cdr(tail)
	}

	return values.NewBool(invariant), nil
}

func LessThanOrImpl(args values.Interface, rt *Runtime, cb Expression) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewBool(true), nil
	}
	head := values.Car(args)
	var (
		lhs, rhs  values.Numeric
		invariant = true
		err       error
	)
	lhs, err = evalToNumber(head, rt, cb)
	if err != nil {
		return values.NewBool(false), err
	}
	for tail := values.Cdr(args); tail.Type() != types.Nil && invariant; tail = values.Cdr(tail) {
		current := values.Car(tail)
		rhs, err = evalToNumber(current, rt, cb)
		if err != nil {
			return values.NewBool(false), err
		}
		invariant = rhs != nil && lhs != nil && lhs.LessThanOrEqual(rhs)
		lhs = rhs
	}

	return values.NewBool(invariant), nil
}

func GreatThanImpl(args values.Interface, rt *Runtime, cb Expression) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewBool(true), nil
	}
	head := values.Car(args)
	var (
		lhs, rhs  values.Numeric
		invariant = true
		err       error
	)
	lhs, err = evalToNumber(head, rt, cb)
	if err != nil {
		return values.NewBool(false), err
	}
	for tail := values.Cdr(args); tail.Type() != types.Nil && invariant; tail = values.Cdr(tail) {
		current := values.Car(tail)
		rhs, err = evalToNumber(current, rt, cb)
		if err != nil {
			return values.NewBool(false), err
		}
		invariant = rhs != nil && lhs != nil && lhs.GreaterThan(rhs)
		lhs = rhs
	}

	return values.NewBool(invariant), nil
}

func GreatThanOrImpl(args values.Interface, rt *Runtime, cb Expression) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewBool(true), nil
	}
	head := values.Car(args)
	var (
		lhs, rhs  values.Numeric
		invariant = true
		err       error
	)
	lhs, err = evalToNumber(head, rt, cb)
	if err != nil {
		return values.NewBool(false), err
	}
	for tail := values.Cdr(args); tail.Type() != types.Nil && invariant; tail = values.Cdr(tail) {
		current := values.Car(tail)
		rhs, err = evalToNumber(current, rt, cb)
		if err != nil {
			return values.NewBool(false), err
		}
		invariant = rhs != nil && lhs != nil && lhs.GreaterThanOrEqual(rhs)
		lhs = rhs
	}

	return values.NewBool(invariant), nil
}

func EqualImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewBool(true), nil
	}
	current := values.Car(args)
	tail := values.Cdr(args)
	for tail.Type() != types.Nil {
		if !current.Equal(values.Car(tail)) {
			return values.NewBool(false), nil
		}
		tail = values.Cdr(tail)
	}
	return values.NewBool(true), nil
}

func QuotImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Pair {
		return values.Car(args), nil
	}
	return args, nil
}

func NotImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args.Type() == types.Nil {
		return values.NewBool(false), ErrWrongNumberOfArguments
	}
	return values.NewBool(!values.Car(args).IsTruthy()), nil
}
