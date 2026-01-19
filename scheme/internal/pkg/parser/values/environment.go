package environment

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"

type Environment struct {
	state map[string]values.Value
}

func NewEnvironment() Environment {
	return Environment{
		state: make(map[string]values.Value),
	}
}

func FromEnvironment(env Environment) Environment {
	return Environment{
		state: env.state,
	}
}

func (env *Environment) Define(name string, value values.Value) {
	env.state[name] = value
}

func (env *Environment) Lookup(name string) (values.Value, bool) {
	v, ok := env.state[name]
	return v, ok
}
