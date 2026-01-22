package parser

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

func (rt *Runtime) defaultEnvironment() {

	rt.Env.Define("newline", values.NewString("\n"))
	//I/O
	rt.Env.Define(types.Format.String(), NewLambda(rt, formatImpl))
	rt.Env.Define(types.Write.String(), NewLambda(rt, writeImpl))
	rt.Env.Define(types.Display.String(), NewLambda(rt, displayImpl))
	//quote
	rt.Env.Define("quot", NewLambda(rt, quotImpl))
	//relational operators
	rt.Env.Define("<", NewLambda(rt, lessThanImpl))
	rt.Env.Define("<=", NewLambda(rt, lessThanOrImpl))
	rt.Env.Define(">", NewLambda(rt, greatThanImpl))
	rt.Env.Define(">=", NewLambda(rt, greatThanOrImpl))
	rt.Env.Define("=", NewLambda(rt, equalImpl))
	//boolean operators
	rt.Env.Define("not", NewLambda(rt, notImpl))
	//arithmetic
	rt.Env.Define("+", NewLambda(rt, sumImpl))
	rt.Env.Define("-", NewLambda(rt, differenceImpl))
	rt.Env.Define("*", NewLambda(rt, productImpl))

}
