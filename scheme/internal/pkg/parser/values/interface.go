package values

import (
	"fmt"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type Interface interface {
	fmt.Stringer
	Equal(p Interface) bool
	Type() types.Type
	SetToken(token lexer.Token)
	GetToken() lexer.Token
	AsPrimitive() (Primitive, error)
	Bool() bool
	DisplayString() string
	WriteString() string
}
