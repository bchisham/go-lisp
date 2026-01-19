package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

type Expression func(args []Value, environment Environment) (Value, error)

type LambdaExpr struct {
	Name string
	Env  Environment
	Body Expression
}

func NewExpression(env Environment, body []Value) Expression {
	return func(args []Value, env Environment) (Value, error) {
		//TODO evaluate body in environment
		return NewVoidType(), nil
	}
}

func NewLambda(env Environment, expression Expression) Value {
	return Value{
		Type: types.Lambda,
		LambdaVal: LambdaExpr{
			Env:  env,
			Body: expression,
		},
	}
}

func (l LambdaExpr) Apply(args []Value) (Value, error) {
	return l.Body(args, l.Env)
}
