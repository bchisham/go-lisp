package builtins

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

func defineImpl(args values.Type, rt *Runtime, evalCallBack Expression) (values.Type, error) {
	switch args.(type) {
	case values.Pair:
		lst := args.(values.Pair)
		head := lst.Car()
		tail := lst.Cdr()
		switch head.(type) {
		case values.Identifier:
			ident := head.(values.Identifier)
			switch tail.(type) {
			case values.Pair:
				vals := tail.(values.Pair)
				if vals.Cdr().Type() == types.Nil {
					rt.Env.Define(ident.GetName(), vals.Car())
					return values.NewVoidType(), nil
				}
			}
			rt.Env.Define(ident.GetName(), tail)
			return values.NewVoidType(), nil
		default:
			return values.NewVoidType(), ErrIllFormedSpecialForm
		}
	default:
		return values.NewVoidType(), ErrWrongNumberOfArguments
	}
}
