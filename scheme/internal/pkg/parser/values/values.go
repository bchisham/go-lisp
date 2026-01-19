package values

import (
	"fmt"
	"strings"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/boolean"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	tokentypes "github.com/bchisham/go-lisp/scheme/internal/pkg/lexer/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type Primitive struct {
	BoolVal   bool
	CharVal   string
	StringVal string
	NameVal   string
	FloatVal  float64
	IntVal    int64
	Literal   string
}

type Value struct {
	Type types.Type
	Primitive
	ListVal   []Value
	LambdaVal LambdaExpr
	Token     lexer.Token
}

func (v *Value) String() string {
	switch v.Type {
	case types.Bool:
		return fmt.Sprintf("%s", boolean.Trinary(v.BoolVal, tokentypes.LiteralTrue, tokentypes.LiteralFalse))
	case types.Char:
		return fmt.Sprintf("%v", v.CharVal)
	case types.String:
		return fmt.Sprintf("%v", v.StringVal)
	case types.Float:
		return fmt.Sprintf("%v", v.FloatVal)
	case types.Int:
		return fmt.Sprintf("%v", v.IntVal)
	case types.List:
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
	case types.Lambda:
		return fmt.Sprintf("%v", "lambda")
	}
	return ""
}

func (v *Value) Equal(other Value) bool {
	if v.Type != other.Type {
		return false
	}
	switch v.Type {
	case types.Bool:
		if v.BoolVal != other.BoolVal {
			return false
		}
	case types.Char:
		if v.CharVal != other.CharVal {
			return false
		}
	case types.String:
		if v.StringVal != other.StringVal {
			return false
		}
	case types.Float:
		if v.FloatVal != other.FloatVal {
			return false
		}
	case types.Int:
		if v.IntVal != other.IntVal {
			return false
		}
	case types.List:
		if len(v.ListVal) != len(other.ListVal) {
			return false
		}
		for i := range v.ListVal {
			if !v.ListVal[i].Equal(other.ListVal[i]) {
				return false
			}
		}
	case types.Identifier, types.RelationalOperator:
		if v.Literal != other.Literal {
			return false
		}
	case types.Lambda:
		return false
	}
	return true
}

func FromToken(tok lexer.Token) (v Value) {
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
		v = NewQuotType()
	case lexer.TokenRelationalOperator:
		v = NewRelationalOperator(tok.Literal)
	case lexer.TokenArithmeticOperator:
		v = NewArithmeticOperator(tok.Literal)
	}
	v.Token = tok
	return
}
