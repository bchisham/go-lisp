package parser

import (
	"errors"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/builtins"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

// evalSexpression evaluate a S-Expression in the given environment
func evalSexpression(l values.Interface, rt *builtins.Runtime) (values.Interface, error) {

	switch l.(type) {
	case values.Nil:
		return values.NewNil(), nil
	case values.Operator:
		lambda, err := applyOperator(l.(values.Operator), rt)
		if err != nil {
			return values.NewVoidType(), err
		}
		return lambda.Apply(values.NewNil())
	case builtins.Lambda:
		return l.(builtins.Lambda).Apply(values.NewNil())
	case values.Quot:
		return values.NewNil(), nil
	case values.Identifier:
		return lookupIdentifier(l.(values.Identifier), rt)
	case values.Pair:
		// continue to S-Expression evaluation
		lst := l.(values.Pair)
		return evaluatePair(lst, rt)
	default:
		return l, ErrOperatorIsNotAProcedure
	}
}

// evaluatePair evaluates a pair as a S-Expression
// It handles the evaluation of the head and applies it to the tail
// according to the rules of the Scheme language
func evaluatePair(lst values.Pair, rt *builtins.Runtime) (values.Interface, error) {
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
			resvVal, err := lookupIdentifier(p, rt)
			if err != nil {
				return values.NewVoidType(), ErrUndefinedIdent
			}
			_, userDefined := types.FromString(p.GetName())
			switch resvVal.(type) {
			case builtins.Lambda:
				if userDefined == nil {
					lambda := resvVal.(builtins.Lambda)
					return lambda.Apply(tail)
				}
				return resvVal, nil
			}
		case values.Operator:
			lambda, err := applyOperator(evaluatedHead.(values.Operator), rt)
			if err != nil {
				return values.NewVoidType(), ErrOperatorIsNotAProcedure
			}
			return lambda.Apply(tail)
		case builtins.Lambda:
			return evaluatedHead.(builtins.Lambda).Apply(tail)
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
		oper, err := applyOperator(head.(values.Operator), rt)
		if err != nil {
			return values.NewVoidType(), ErrOperatorIsNotAProcedure
		}
		return oper.Apply(tail)
	case builtins.Lambda:
		return head.(builtins.LambdaExpr), nil
	case values.Identifier:
		p := head.(values.Identifier)
		resolvVal, err := lookupIdentifier(p, rt)
		if err != nil {
			return values.NewVoidType(), ErrUndefinedIdent
		}
		switch resolvVal.(type) {
		case builtins.Lambda:
			return resolvVal.(builtins.Lambda).Apply(tail)
		default:
			return resolvVal, nil
		}
	default:
		if tail.Type() == types.Nil {
			return head, nil
		}
	}
	return evalSexpression(values.Cons(tail, head), rt)
}

func applyOperator(p values.Operator, rt *builtins.Runtime) (builtins.Lambda, error) {

	oper, ok := rt.Env.Lookup(p.GetName())
	if !ok {
		return nil, ErrOperatorIsNotAProcedure
	}
	lambda, ok := oper.(builtins.LambdaExpr)
	if !ok {
		return nil, ErrOperatorIsNotAProcedure
	}
	return lambda, nil
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
