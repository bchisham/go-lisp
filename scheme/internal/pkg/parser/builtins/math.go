package builtins

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

var ArithmeticAllowedGate = types.NewTypeGate(types.Int, types.Float)

// evalToNumber evaluates val using the evaluationCallback if it is not already a numeric type.
// It returns the evaluated numeric value or an error if the value is not numeric.
func evalToNumber(val values.Interface, rt *Runtime, evaluationCallback Expression) (values.Numeric, error) {
	if ArithmeticAllowedGate(val.Type()) {
		num, ok := val.(values.Numeric)
		if !ok {
			return values.Zero, ErrNumberExpected
		}
		return num, nil
	}
	evaluated, err := evaluationCallback(val, rt)
	if err != nil {
		return values.Zero, err
	}
	if !ArithmeticAllowedGate(evaluated.Type()) {
		return values.Zero, ErrNumberExpected
	}
	num, ok := evaluated.(values.Numeric)
	if !ok {
		return values.Zero, ErrNumberExpected
	}
	return num, nil

}

// SumImpl computes the sum of all numeric arguments in args.
func SumImpl(args values.Interface, rt *Runtime, evaluationCallback Expression) (_ values.Interface, err error) {
	if args.Type() == types.Nil {
		return values.NewInt(0), nil
	}
	if ArithmeticAllowedGate(args.Type()) {
		return args, nil
	}
	if ArithmeticAllowedGate(values.Car(args).Type()) && values.Cdr(args).Type() == types.Nil {
		return values.Car(args), nil
	}

	sum := values.Zero
	for args.Type() != types.Nil {
		head := values.Car(args)
		lhs, err := evalToNumber(head, rt, evaluationCallback)
		if err != nil {
			return values.NewNil(), err
		}
		sum = sum.Add(lhs)
		args = values.Cdr(args)
	}
	return sum, nil
}

// DifferenceImpl computes the difference of all numeric arguments in args.
// It subtracts each subsequent number from the first.
func DifferenceImpl(args values.Interface, rt *Runtime, evaluationCallback Expression) (_ values.Interface, err error) {
	if args.Type() == types.Nil {
		return values.NewInt(0), nil
	}
	if ArithmeticAllowedGate(args.Type()) {
		return args, nil
	}
	diff := values.Zero
	for args.Type() != types.Nil {
		current := values.Car(args)
		lhs, err := evalToNumber(current, rt, evaluationCallback)
		if err != nil {
			return values.NewNil(), err
		}
		diff = diff.Sub(lhs)
		args = values.Cdr(args)
	}

	return diff, nil
}

// ProductImpl computes the product of all numeric arguments in args.
// It multiplies each number together.
func ProductImpl(args values.Interface, rt *Runtime, evaluationCallback Expression) (_ values.Interface, err error) {
	if args.Type() == types.Nil {
		return values.NewInt(1), nil
	}

	if ArithmeticAllowedGate(args.Type()) {
		return args, nil
	}

	product := values.One
	for args.Type() != types.Nil {
		current := values.Car(args)
		lhs, err := evalToNumber(current, rt, evaluationCallback)
		if err != nil {
			return values.NewNil(), err
		}
		product = product.Mul(lhs)
		args = values.Cdr(args)
	}
	return product, nil
}

// QuotientImpl computes the quotient of all numeric arguments in args.
// It divides the first number by each subsequent number in order.
func QuotientImpl(args values.Interface, rt *Runtime, evaluationCallback Expression) (_ values.Interface, err error) {
	if args.Type() == types.Nil {
		return values.NewNil(), ErrWrongNumberOfArguments
	}
	if ArithmeticAllowedGate(args.Type()) {
		return args, nil
	}
	tail := values.Cdr(args)

	if tail.Type() == types.Nil {
		return values.NewNil(), ErrWrongNumberOfArguments
	}
	var quotient = values.One
	for tail.Type() != types.Nil {
		current := values.Car(tail)
		rhs, err := evalToNumber(current, rt, evaluationCallback)
		if err != nil {
			return values.NewNil(), err
		}
		if rhs.Equal(values.Zero) {
			return values.NewNil(), ErrDivideByZero
		}
		quotient, err = quotient.Div(rhs)
		if err != nil {
			return values.NewNil(), err
		}
		tail = values.Cdr(tail)
	}

	return quotient, nil
}

// RemainderImpl computes the remainder of the division of the first numeric argument by each subsequent numeric argument in args.
func RemainderImpl(args values.Interface, rt *Runtime, evaluationCallback Expression) (_ values.Interface, err error) {
	if args.Type() == types.Nil {
		return values.NewNil(), ErrWrongNumberOfArguments
	}

	if ArithmeticAllowedGate(args.Type()) {
		return args, nil
	}
	current := values.Car(args)
	tail := values.Cdr(args)
	var remainder = values.Zero
	for tail.Type() != types.Nil {

		rhsNum, err := evalToNumber(current, rt, evaluationCallback)
		if err != nil {
			return values.NewNil(), ErrNumberExpected
		}
		if rhsNum.Equal(values.Zero) {
			return values.NewNil(), ErrDivideByZero
		}
		remainder, err = remainder.Mod(rhsNum)
		if err != nil {
			return values.NewNil(), err
		}
		current = values.Car(tail)
		tail = values.Cdr(tail)
	}

	return remainder, nil
}
