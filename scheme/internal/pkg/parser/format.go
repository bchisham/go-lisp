package parser

import (
	"errors"
	"fmt"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

var (
	ErrInvalidFormat = errors.New("invalid format")
	ErrBadArgument   = errors.New("bad argument")
)

func formatImpl(args []types.Value) (types.Value, error) {
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
