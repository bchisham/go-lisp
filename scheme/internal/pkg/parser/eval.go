package parser

import (
	"errors"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/builtins"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

// evalSexpression evaluate a S-Expression i the given environment.
func evalSexpression(l values.Interface, rt *builtins.Runtime) (values.Interface, error) {

	switch l.(type) {
	case values.Nil:
		return values.NewNil(), nil
	case values.Operator:
		p, ok := l.(values.Operator)
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		oper, ok := rt.Env.Lookup(p.GetName())
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		lambda, ok := oper.(builtins.LambdaExpr)
		if !ok {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		return lambda.Apply(values.NewNil())
	case builtins.LambdaExpr:
		lambda := l.(builtins.LambdaExpr)

		return lambda.Apply(values.NewNil())
	case values.Quot:
		return values.NewNil(), nil
	case values.Identifier:
		p, ok := l.(values.Identifier)
		resolvVal, ok := rt.Env.Lookup(p.GetName())
		if !ok {
			return values.NewVoidType(), ErrUndefinedIdent
		}
		return resolvVal, nil
	case values.Pair:
		// continue to S-Expression evaluation
		lst := l.(values.Pair)
		head := lst.Car()
		tail := lst.Cdr()
		switch head.(type) {
		case values.Pair:
			evaluatedHead, err := evalSexpression(head, rt)
			if err != nil && !errors.Is(err, ErrOperatorIsNotAProcedure) {
				return evaluatedHead, err
			}
			switch evaluatedHead.(type) {
			case values.Identifier:
				p := evaluatedHead.(values.Identifier)
				resvVal, ok := rt.Env.Lookup(p.GetName())
				if !ok {
					return values.NewVoidType(), ErrUndefinedIdent
				}
				_, userDefined := types.FromString(p.GetName())
				switch resvVal.(type) {
				case builtins.LambdaExpr:
					if userDefined == nil {
						lambda := resvVal.(builtins.LambdaExpr)
						return lambda.Apply(tail)
					}
					return resvVal, nil
				}
			case values.Operator:
				p := evaluatedHead.(values.Operator)
				oper, ok := rt.Env.Lookup(p.GetName())
				if !ok {
					return values.NewVoidType(), ErrOperatorIsNotAProcedure
				}
				lambda, ok := oper.(builtins.Lambda)
				return lambda.Apply(tail)
			case builtins.Lambda:
				lambda := evaluatedHead.(builtins.Lambda)
				return lambda.Apply(tail)
			case values.Quot:
				return tail, nil
			case values.Numeric, values.Boolean, values.String, values.Char:
				if tail.Type() == types.Nil {
					return evaluatedHead, nil
				}
				return evalSexpression(values.Cons(tail, evaluatedHead), rt)
			default:
				return evaluatedHead, ErrOperatorIsNotAProcedure
			}
		case values.Quot:
			switch tail.Type() {
			case types.Nil:
				return values.NewNil(), ErrWrongNumberOfArguments
			default:
				return values.NewQuot(tail), nil
			}
		case values.Operator:
			p := head.(values.Operator)
			oper, ok := rt.Env.Lookup(p.GetName())
			if !ok {
				return values.NewVoidType(), ErrOperatorIsNotAProcedure
			}
			lambda, ok := oper.(builtins.LambdaExpr)
			if !ok {
				return values.NewVoidType(), ErrOperatorIsNotAProcedure
			}
			return lambda.Apply(tail)
		case builtins.Lambda:
			lambda := head.(builtins.LambdaExpr)
			return lambda, nil
		case values.Identifier:
			p := head.(values.Identifier)
			resolvVal, ok := rt.Env.Lookup(p.GetName())
			if !ok {
				return values.NewVoidType(), ErrUndefinedIdent
			}
			switch resolvVal.(type) {
			case builtins.Lambda:
				lambda := resolvVal.(builtins.Lambda)
				return lambda.Apply(tail)
			default:
				return resolvVal, nil
			}
		default:
			if tail.Type() == types.Nil {
				return head, nil
			}
		}
	default:
		return l, ErrOperatorIsNotAProcedure
	}
	return l, ErrOperatorIsNotAProcedure
}

func lookupProcedure(ident values.Identifier, rt *builtins.Runtime) (builtins.Lambda, error) {
	resolvVal, ok := rt.Env.Lookup(ident.GetName())
	if !ok {
		return nil, ErrUndefinedIdent
	}
	lambda, ok := resolvVal.(builtins.Lambda)
	if !ok {
		return nil, ErrOperatorIsNotAProcedure
	}
	return lambda, nil
}

func lookupIdentifier(ident values.Identifier, rt *builtins.Runtime) (values.Interface, error) {
	resolvVal, ok := rt.Env.Lookup(ident.GetName())
	if !ok {
		return values.NewVoidType(), ErrUndefinedIdent
	}
	return resolvVal, nil
}
