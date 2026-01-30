package builtins

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

type Environment struct {
	state map[string]values.Interface
}

func NewEnvironment() Environment {
	return Environment{
		state: make(map[string]values.Interface),
	}
}

func FromEnvironment(env Environment) Environment {
	return Environment{
		state: env.state,
	}
}

func (env *Environment) Define(name string, value values.Interface) {
	env.state[name] = value
}

func (env *Environment) Lookup(name string) (values.Interface, bool) {
	v, ok := env.state[name]
	return v, ok
}

func adaptBuiltin(f func(value values.Interface, rt *Runtime, cb Expression) (values.Interface, error), cb Expression) Expression {
	return func(args values.Interface, rt *Runtime) (values.Interface, error) {
		return f(args, rt, cb)
	}
}

func (rt *Runtime) defaultEnvironment(cb Expression) {

	rt.Env.Define("newline", values.NewString("\n"))
	//I/O
	rt.Env.Define(types.Format.String(), NewLambda(rt, FormatImpl))
	rt.Env.Define(types.Write.String(), NewLambda(rt, WriteImpl))
	rt.Env.Define(types.Display.String(), NewLambda(rt, DisplayImpl))
	//quote
	rt.Env.Define("quot", NewLambda(rt, QuotImpl))
	//relational operators
	rt.Env.Define("<", NewLambda(rt, adaptBuiltin(LessThanImpl, cb)))
	rt.Env.Define("<=", NewLambda(rt, adaptBuiltin(LessThanOrImpl, cb)))
	rt.Env.Define(">", NewLambda(rt, adaptBuiltin(GreatThanImpl, cb)))
	rt.Env.Define(">=", NewLambda(rt, adaptBuiltin(GreatThanOrImpl, cb)))
	rt.Env.Define("=", NewLambda(rt, EqualImpl))
	//boolean operators
	rt.Env.Define("not", NewLambda(rt, NotImpl))
	//arithmetic
	rt.Env.Define("+", NewLambda(rt, adaptBuiltin(SumImpl, cb)))
	rt.Env.Define("-", NewLambda(rt, adaptBuiltin(DifferenceImpl, cb)))
	rt.Env.Define("*", NewLambda(rt, adaptBuiltin(ProductImpl, cb)))
	rt.Env.Define("/", NewLambda(rt, adaptBuiltin(QuotientImpl, cb)))
	rt.Env.Define("modulo", NewLambda(rt, adaptBuiltin(RemainderImpl, cb)))

}
