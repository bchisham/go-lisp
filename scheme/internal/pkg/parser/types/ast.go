package types

import (
	"fmt"
	"strings"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/boolean"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
)

type Type string

const (
	Bool       Type = "bool"
	Char       Type = "char"
	Float      Type = "float"
	Int        Type = "int"
	List       Type = "list"
	Lambda     Type = "lambda"
	Map        Type = "map"
	String     Type = "string"
	Identifier Type = "identifier"
	Void       Type = "void"
)

type Primitive struct {
	BoolVal   bool
	CharVal   string
	StringVal string
	NameVal   string
	FloatVal  float64
	IntVal    int64
}

type Expression func(args []Value, environment Environment) (Value, error)
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

type LambdaExpr struct {
	Name string
	Env  Environment
	Body Expression
}

type Value struct {
	Type Type
	Primitive
	ListVal   []Value
	LambdaVal LambdaExpr
}

func (v *Value) String() string {
	switch v.Type {
	case Bool:
		return fmt.Sprintf("%s", boolean.Trinary(v.BoolVal, types.LiteralTrue, types.LiteralFalse))
	case Char:
		return fmt.Sprintf("%v", v.CharVal)
	case String:
		return fmt.Sprintf("%v", v.StringVal)
	case Float:
		return fmt.Sprintf("%v", v.FloatVal)
	case Int:
		return fmt.Sprintf("%v", v.IntVal)
	case List:
		sb := strings.Builder{}
		sb.WriteString("(")
		sb.WriteString(
			strings.Join(
				list.Apply(
					v.ListVal,
					func(v Value) string { return v.String() }),
				" "),
		)
		sb.WriteString(")")
		return fmt.Sprintf("%s", sb.String())
	case Lambda:
		return fmt.Sprintf("%v", "lambda")
	}
	return ""
}

func (l LambdaExpr) Apply(args []Value) (Value, error) {
	return l.Body(args, l.Env)
}
