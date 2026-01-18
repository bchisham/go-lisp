package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

var (
	ErrInvalidFormat = errors.New("invalid format")
	ErrBadArgument   = errors.New("bad argument")
)

func displayImpl(args []types.Value, env types.Environment) (types.Value, error) {
	sb := &strings.Builder{}
	sb.WriteString(strings.Join(list.Apply(args, func(v types.Value) string { return v.String() }), " "))
	fmt.Println(sb.String())
	return newVoidType(), nil
}

func writeImpl(args []types.Value, env types.Environment) (types.Value, error) {
	sb := &strings.Builder{}
	sb.WriteString(strings.Join(list.Apply(args, func(v types.Value) string {
		switch v.Type {
		case types.String:
			return "\"" + v.String() + "\""
		default:
			return v.String()
		}
	}), " "))
	fmt.Print(sb.String())
	return newVoidType(), nil

}

func formatImpl(args []types.Value, env types.Environment) (types.Value, error) {
	if len(args) < 2 {
		return newVoidType(), ErrBadArgument
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

	return newVoidType(), nil
}
