package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

func NewVoidType() Interface {
	return void{}
}

type void struct {
	truthyValue
}

func (v void) Equal(p Interface) bool {
	_, ok := p.(void)
	return ok
}

func (v void) Type() types.Type {
	return types.Void
}

func (v void) DisplayString() string {
	return "#<void>"
}

func (v void) WriteString() string {
	return "#<void>"
}
