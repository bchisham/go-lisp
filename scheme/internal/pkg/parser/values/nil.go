package values

import (
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type Nil struct {
	truthyValue
}

func NewNil() Interface {
	return Nil{}
}

func (n Nil) String() string {
	return "()"
}

func (n Nil) Equal(p Interface) bool {
	_, ok := p.(Nil)
	return ok
}

func (n Nil) Type() types.Type {
	return types.Nil
}

func (n Nil) DisplayString() string {
	return "()"
}

func (n Nil) WriteString() string {
	return "()"
}
