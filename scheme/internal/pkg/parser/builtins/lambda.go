package builtins

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

type Lambda interface {
	values.Type
	Apply(args values.Type) (values.Type, error)
}

type Expression func(args values.Type, rt *Runtime) (values.Type, error)

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

func (l LambdaExpr) Equal(p values.Type) bool {

	return false
}

func (l LambdaExpr) Type() types.Type {
	return types.Lambda
}

func (l LambdaExpr) SetToken(token lexer.Token) {
	l.srcToken = token
}

func NewExpression(env Environment, body values.Type) Expression {
	return func(args values.Type, rt *Runtime) (values.Type, error) {
		//TODO evaluate body in environment
		return values.NewVoidType(), nil
	}
}

func NewLambda(rt *Runtime, expression Expression) values.Type {
	return LambdaExpr{
		Runtime: rt,
		Body:    expression,
	}
}

func (l LambdaExpr) Apply(args values.Type) (values.Type, error) {
	return l.Body(args, l.Runtime)
}

func (l LambdaExpr) IsTruthy() bool {
	return true
}

func lambdaImpl(args values.Type, rt *Runtime, evalCallBack Expression) (values.Type, error) {
	return NewLambda(rt, NewExpression(rt.Env, args)), nil
}
