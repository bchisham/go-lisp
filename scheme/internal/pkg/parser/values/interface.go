package values

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type Type interface {
	Equal(p Type) bool
	Type() types.Type
	IsTruthy() bool
	DisplayString() string
	WriteString() string
}

type Numeric interface {
	Type
	IsInteger() bool
	IsFloat() bool
	AsFloat() (float64, error)
	AsInt() (int64, error)
	Add(rhs Numeric) Numeric
	Sub(rhs Numeric) Numeric
	Mul(rhs Numeric) Numeric
	Div(rhs Numeric) (Numeric, error)
	Mod(rhs Numeric) (Numeric, error)
	LessThan(p Numeric) bool
	LessThanOrEqual(p Numeric) bool
	GreaterThan(p Numeric) bool
	GreaterThanOrEqual(p Numeric) bool
}
