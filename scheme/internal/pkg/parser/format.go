package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

var (
	ErrInvalidFormat = errors.New("invalid format")
	ErrBadArgument   = errors.New("bad argument")
)

func displayImpl(args []values.Value, env values.Environment) (values.Value, error) {
	sb := &strings.Builder{}
	sb.WriteString(strings.Join(list.Apply(args, func(v values.Value) string { return v.String() }), " "))
	fmt.Println(sb.String())
	return values.NewVoidType(), nil
}

func writeImpl(args []values.Value, env values.Environment) (values.Value, error) {
	sb := &strings.Builder{}
	sb.WriteString(strings.Join(list.Apply(args, func(v values.Value) string {
		switch v.Type {
		case types.String:
			switch v.String() {
			//case "\n":
			//	return "\n"
			default:
				return fmt.Sprintf("%q", v.String())
			}
		default:
			return v.String()
		}
	}), " "))
	fmt.Print(sb.String())
	return values.NewVoidType(), nil

}

func formatImpl(args []values.Value, env values.Environment) (values.Value, error) {
	if len(args) < 2 {
		return values.NewVoidType(), ErrBadArgument
	}
	f := list.Car(args)
	obj := list.Car(list.Cdr(args))

	switch f.NameVal {
	case "t":
		switch obj.Type {
		case types.String:
			fmt.Printf("%s", obj.StringVal)
		case types.Identifier:
			fmt.Printf("%s", obj.NameVal)
		}
	default:
		fmt.Printf("%#v", obj)
	}

	return values.NewVoidType(), nil
}
