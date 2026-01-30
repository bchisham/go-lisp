package builtins

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

type Lambda interface {
	values.Interface
	Apply(args values.Interface) (values.Interface, error)
}

type Expression func(args values.Interface, rt *Runtime) (values.Interface, error)

type LambdaExpr struct {
	Name     string
	Runtime  *Runtime
	Body     Expression
	srcToken lexer.Token
}

func (l LambdaExpr) WriteString() string {
	return "#<procedure>"
}

func (l LambdaExpr) DisplayString() string {
	return "#<procedure>"
}

func (l LambdaExpr) String() string {
	return "#<procedure>"

}

func (l LambdaExpr) Equal(p values.Interface) bool {

	return false
}

func (l LambdaExpr) Type() types.Type {
	return types.Lambda
}

func (l LambdaExpr) SetToken(token lexer.Token) {
	l.srcToken = token
}

func NewExpression(env Environment, body []values.Interface) Expression {
	return func(args values.Interface, rt *Runtime) (values.Interface, error) {
		//TODO evaluate body in environment
		return values.NewVoidType(), nil
	}
}

func NewLambda(rt *Runtime, expression Expression) values.Interface {
	return LambdaExpr{
		Runtime: rt,
		Body:    expression,
	}
}

func (l LambdaExpr) Apply(args values.Interface) (values.Interface, error) {
	return l.Body(args, l.Runtime)
}

func (l LambdaExpr) IsTruthy() bool {
	return true
}
