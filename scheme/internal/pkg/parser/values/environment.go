package values

type Environment struct {
	state map[string]Value
}

func NewEnvironment() Environment {
	return Environment{
		state: make(map[string]Value),
	}
}

func FromEnvironment(env Environment) Environment {
	return Environment{
		state: env.state,
	}
}

func (env *Environment) Define(name string, value Value) {
	env.state[name] = value
}

func (env *Environment) Lookup(name string) (Value, bool) {
	v, ok := env.state[name]
	return v, ok
}
