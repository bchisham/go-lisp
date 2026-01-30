package values

import (
	"errors"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
)

var ErrNotAPrimitive = errors.New("not a primitive")

func FromToken(tok lexer.Token) (v Interface) {

	switch tok.Type {
	case lexer.TokenIdent:
		v = NewIdentifier(tok.Literal)
	case lexer.TokenInt:
		v = NewInt(tok.Int)
	case lexer.TokenBoolean:
		v = NewBool(tok.Bool)
	case lexer.TokenString:
		v = NewString(tok.Literal)
	case lexer.TokenQuot:
		v = NewQuot(NewNil())
	case lexer.TokenRelationalOperator:
		v = NewRelationalOperator(tok.Literal)
	case lexer.TokenArithmeticOperator:
		v = NewArithmeticOperator(tok.Literal)
	}

	return
}
