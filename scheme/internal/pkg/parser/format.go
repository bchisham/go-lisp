package parser

import (
	"fmt"
	"strings"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

var ()

func displayImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	sb := &strings.Builder{}
	sb.WriteString(strings.Join(list.Apply(args, func(v values.Interface) string { return v.DisplayString() }), " "))
	_, err := fmt.Fprintln(rt.Out, sb.String())
	if err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}

func writeImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	sb := &strings.Builder{}
	sb.WriteString(strings.Join(list.Apply(args, func(v values.Interface) string {
		switch v.Type() {
		case types.String:
			switch v.String() {
			default:
				return fmt.Sprintf("%q", v.WriteString())
			}
		default:
			return v.WriteString()
		}
	}), " "))
	_, err := fmt.Fprint(rt.Out, sb.String())
	if err != nil {
		return values.NewVoidType(), ErrIo(err)
	}
	return values.NewVoidType(), nil
}

func formatImpl(args []values.Interface, rt *Runtime) (values.Interface, error) {
	if len(args) < 2 {
		return values.NewVoidType(), ErrBadArgument
	}
	f := list.Car(args)
	obj := list.Car(list.Cdr(args))
	var err error
	if f.Type() != types.Identifier {
		return values.NewVoidType(), ErrBadArgument
	}

	switch f.String() {
	case "t":
		switch obj.Type() {
		case types.String:
			_, err = fmt.Fprintf(rt.Out, "%s", obj.String())
			return values.NewVoidType(), ErrIo(err)
		case types.Identifier:
			_, err = fmt.Fprintf(rt.Out, "%s", obj.String())
			return values.NewVoidType(), ErrIo(err)
		}
	default:
		_, err = fmt.Fprintf(rt.Out, "%#v", obj)
		if err != nil {
			return values.NewVoidType(), ErrIo(err)
		}
	}

	return values.NewVoidType(), nil
}
