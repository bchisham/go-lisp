package builtins

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

func setBangImpl(args values.Type, rt *Runtime, evalCallBack Expression) (values.Type, error) {
	panic("implement me")
}

// letImpl implements the 'let' special form.
// It evaluates a series of bindings and then evaluates the body expressions
// in an environment extended with those bindings.
// The 'let' form has the syntax:
// (let ((var1 val1) (var2 val2) ...) body1 body2 ...)
// It returns the result of evaluating the body expressions.
func letImpl(args values.Type, rt *Runtime, evalCallBack Expression) (values.Type, error) {

	switch args.(type) {
	case values.Pair:
		lst := args.(values.Pair)
		bindings := lst.Car()
		body := lst.Cdr()

		// Create a new environment extended from the current one
		childRt := rt.ExtendEnvironment()
		// Process bindings
		switch bindings.(type) {
		case values.Nil:
			// No bindings, proceed to evaluate body
		case values.Pair:
			bindList := bindings
			for bindList.Type() != types.Nil {
				binding := values.Car(bindList)
				switch binding.(type) {
				case values.Pair:
					bindPair := binding.(values.Pair)
					varNameType := bindPair.Car()
					varValueType := bindPair.Cdr()
					// Evaluate the value expression
					evaluatedValue, err := evalCallBack(varValueType, rt)
					if err != nil {
						return values.NewVoidType(), err
					}
					// Define the variable in the new environment
					switch varNameType.(type) {
					case values.Identifier:
						varName := varNameType.(values.Identifier).GetName()
						childRt.Env.Define(varName, evaluatedValue)
					default:
						return values.NewVoidType(), ErrIllFormedSpecialForm
					}
				default:
					return values.NewVoidType(), ErrIllFormedSpecialForm
				}
				bindList = values.Cdr(bindList)
			}
		default:
			return values.NewVoidType(), ErrIllFormedSpecialForm
		}

		// Evaluate body expressions in the new environment
		var result values.Type = values.NewVoidType()
		for body.Type() != types.Nil {
			currentExpr := values.Car(body)
			resultEval, err := evalCallBack(currentExpr, childRt)
			if err != nil {
				return values.NewVoidType(), err
			}
			result = resultEval
			body = values.Cdr(body)
		}
		return result, nil
	default:
		return values.NewVoidType(), ErrWrongNumberOfArguments
	}

}
