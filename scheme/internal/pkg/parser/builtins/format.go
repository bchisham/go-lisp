package builtins

import (
	"fmt"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

func DisplayImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args == nil {
		return values.NewVoidType(), ErrBadArgument
	}
	//TODO if cdr is not nil handle PORT implementation

	// Print the display string of the value
	if _, err := fmt.Fprintf(rt.Out, "%s", args.DisplayString()); err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}

func WriteImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if _, err := fmt.Fprintf(rt.Out, "%s", args.WriteString()); err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}

func FormatImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if _, err := fmt.Fprintf(rt.Out, "%s", args.DisplayString()); err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}
