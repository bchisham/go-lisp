package values

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type Interface interface {
	Equal(p Interface) bool
	Type() types.Type
	AsPrimitive() (Primitive, error)
	IsTruthy() bool
	DisplayString() string
	WriteString() string
}
