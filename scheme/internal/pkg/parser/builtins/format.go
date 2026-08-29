package builtins

import (
	"errors"
	"fmt"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

func DisplayImpl(args values.Type, rt *Runtime) (values.Type, error) {
	if args == nil {
		return values.NewVoidType(), ErrBadArgument
	}

	// Print the display string of the value
	evalArgs, err := rt.Eval(args)
	if err != nil {
		if errors.Is(err, ErrOperatorIsNotAProcedure) {
			// If the error is that the operator is not a procedure, we want to print the string representation of the arguments instead of returning an error
			if _, err := fmt.Fprintf(rt.Out, "%s", args.DisplayString()); err != nil {
				return values.NewVoidType(), ErrIo(err)
			}
			return values.NewVoidType(), nil
		}
		return values.NewVoidType(), err
	}

	if _, err := fmt.Fprintf(rt.Out, "%s", evalArgs.DisplayString()); err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}

func WriteImpl(args values.Type, rt *Runtime) (values.Type, error) {

	evalArgs, err := rt.Eval(args)
	if err != nil {
		// If the error is that the operator is not a procedure, we want to print the string representation of the arguments instead of returning an error
		if errors.Is(err, ErrOperatorIsNotAProcedure) {
			if _, err := fmt.Fprintf(rt.Out, "%s", args.WriteString()); err != nil {
				return values.NewVoidType(), ErrIo(err)
			}
			return values.NewVoidType(), nil
		}
		return values.NewVoidType(), err
	}
	if _, err := fmt.Fprintf(rt.Out, "%s", evalArgs.WriteString()); err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}

func FormatImpl(args values.Type, rt *Runtime) (values.Type, error) {
	evalArgs, err := rt.Eval(args)
	if err != nil {
		if errors.Is(err, ErrOperatorIsNotAProcedure) {
			if _, err := fmt.Fprintf(rt.Out, "%s", args.DisplayString()); err != nil {
				return values.NewVoidType(), ErrIo(err)
			}
			return values.NewVoidType(), nil
		}
		return values.NewVoidType(), err
	}
	if _, err := fmt.Fprintf(rt.Out, "%s", evalArgs.DisplayString()); err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}
