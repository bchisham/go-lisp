package builtins

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

// Environment represents the runtime environment holding variable bindings.
// It maps variable names to their corresponding values.
// It supports defining new variables and looking up existing ones.
type Environment struct {
	state map[string]values.Interface
}

func NewEnvironment() Environment {
	return Environment{
		state: make(map[string]values.Interface),
	}
}

// ExtendEnvironment creates a new Environment by copying the state from an existing one.
// This allows for creating a new environment that shares the same variable bindings as the original.
func ExtendEnvironment(env Environment) Environment {
	return Environment{
		state: env.state,
	}
}

// Define adds a new variable binding to the environment.
// It associates the given name with the provided value.
func (env *Environment) Define(name string, value values.Interface) {
	env.state[name] = value
}

// Lookup retrieves the value associated with the given variable name.
// It returns the value and a boolean indicating whether the variable was found.
// If the variable is not found, the boolean will be false.
func (env *Environment) Lookup(name string) (values.Interface, bool) {
	v, ok := env.state[name]
	return v, ok
}

// adaptBuiltin adapts a built-in function to match the expected Expression signature.
// It takes a function 'f' that accepts a values.Interface, a Runtime pointer, and an Expression callback,
// and returns an Expression that only requires values.Interface and Runtime pointer.
// This allows built-in functions to be used seamlessly within the expression evaluation framework.
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
