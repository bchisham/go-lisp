package parser

import (
	"fmt"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

func displayImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if args == nil {
		return values.NewVoidType(), ErrBadArgument
	}
	val := values.NewNil()
	if args.Type() == types.Pair {
		val = values.Car(args)
	}
	//TODO if cdr is not nil handle PORT implementation

	// Print the display string of the value
	if _, err := fmt.Fprintf(rt.Out, "%s", val.DisplayString()); err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}

func writeImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if _, err := fmt.Fprintf(rt.Out, "%s", args.WriteString()); err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}

func formatImpl(args values.Interface, rt *Runtime) (values.Interface, error) {
	if _, err := fmt.Fprintf(rt.Out, "%s", args.DisplayString()); err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}
