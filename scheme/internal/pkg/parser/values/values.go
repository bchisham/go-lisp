package values

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/boolean"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	tokentypes "github.com/bchisham/go-lisp/scheme/internal/pkg/lexer/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"golang.org/x/exp/slices"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

var ErrNotAPrimitive = errors.New("not a primitive")

type Primitive struct {
	BoolVal   bool
	CharVal   rune
	StringVal string
	NameVal   string
	FloatVal  float64
	IntVal    int64
	Literal   string
}

type Value struct {
	t types.Type
	Primitive
	ListVal []Interface
	Token   lexer.Token
}

func (v Value) Type() types.Type {
	return v.t
}

func (v Value) AsPrimitive() (Primitive, error) {
	if slices.Contains(list.New(
		types.Bool,
		types.Int,
		types.Float,
		types.Char,
		types.String,
		types.Identifier,
		types.RelationalOperator,
		types.ArithmeticOperator), v.t) {
		return v.Primitive, nil
	}
	return Primitive{}, ErrNotAPrimitive
}

func (v Value) String() string {
	return v.WriteString()
}

func (v Value) DisplayString() string {
	switch v.Type() {
	case types.Bool:
		// Assuming tokentypes.LiteralTrue is "#t" etc.
		return boolean.Trinary(v.BoolVal, tokentypes.LiteralTrue, tokentypes.LiteralFalse)

	case types.Char:
		// display prints the character itself
		return string(v.CharVal)

	case types.String:
		// display prints raw string contents
		return v.StringVal

	case types.Float:
		return strconv.FormatFloat(v.FloatVal, 'g', -1, 64)

	case types.Int:
		return strconv.FormatInt(v.IntVal, 10)

	case types.List:
		var sb strings.Builder
		sb.WriteByte('(')
		elems := list.Apply(v.ListVal, func(x Interface) string {
			// Important: decide what Interface is. If it's your Value,
			// call DisplayString; if it's fmt.Stringer, this will call WriteString via String().
			// Prefer: assert to Value.
			if vv, ok := x.(Value); ok {
				return vv.DisplayString()
			}
			return fmt.Sprint(x)
		})
		sb.WriteString(strings.Join(elems, " "))
		sb.WriteByte(')')
		return sb.String()

	case types.Lambda:
		// display usually prints something like #<procedure>
		return "#<procedure>"

	case types.Void:
		// display of void normally prints nothing; REPL decides whether to show anything.
		return ""
	}
	return ""
}

func (v Value) WriteString() string {
	switch v.Type() {
	case types.Bool:
		return boolean.Trinary(v.BoolVal, tokentypes.LiteralTrue, tokentypes.LiteralFalse)

	case types.Char:
		// write prints a readable char literal.
		// Pick a convention; R5RS-ish is #\a, #\space, #\newline, etc.
		return writeChar(v.CharVal)

	case types.String:
		// write prints a quoted, escaped string
		return strconv.Quote(v.StringVal)

	case types.Float:
		return strconv.FormatFloat(v.FloatVal, 'g', -1, 64)

	case types.Int:
		return strconv.FormatInt(v.IntVal, 10)

	case types.List:
		var sb strings.Builder
		sb.WriteByte('(')
		elems := list.Apply(v.ListVal, func(x Interface) string {
			if vv, ok := x.(Value); ok {
				return vv.WriteString()
			}
			return fmt.Sprint(x)
		})
		sb.WriteString(strings.Join(elems, " "))
		sb.WriteByte(')')
		return sb.String()

	case types.Lambda:
		return "#<procedure>"

	case types.Void:
		// write of void is typically unspecified; keep it empty and let REPL hide it.
		return ""
	}
	return ""
}

func (v Value) Equal(other Interface) bool {
	if v.Type() != other.Type() {
		return false
	}
	o, ok := other.(Value)
	if !ok {
		return false
	}
	switch v.Type() {
	case types.Bool:
		if v.BoolVal != o.BoolVal {
			return false
		}
	case types.Char:
		if v.CharVal != o.CharVal {
			return false
		}
	case types.String:
		if v.StringVal != o.StringVal {
			return false
		}
	case types.Float:
		if v.FloatVal != o.FloatVal {
			return false
		}
	case types.Int:
		if v.IntVal != o.IntVal {
			return false
		}
	case types.List:
		if len(v.ListVal) != len(o.ListVal) {
			return false
		}
		for i := range v.ListVal {
			if !v.ListVal[i].Equal(o.ListVal[i]) {
				return false
			}
		}
	case types.Identifier, types.RelationalOperator:
		if v.Literal != o.Literal {
			return false
		}
	case types.Lambda:
		return false
	}
	return true
}

func (v Value) SetToken(token lexer.Token) {
	v.Token = token
}

func (v Value) GetToken() lexer.Token {
	return v.Token
}

func (v Value) Bool() bool {
	switch v.Type() {
	case types.Bool:
		return v.BoolVal
	default:
		return true
	}
}

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
		v = NewQuotType()
	case lexer.TokenRelationalOperator:
		v = NewRelationalOperator(tok.Literal)
	case lexer.TokenArithmeticOperator:
		v = NewArithmeticOperator(tok.Literal)
	}
	v.SetToken(tok)
	return
}

func FromPrimitive(t types.Type, p Primitive) (v Interface) {
	return Value{
		t:         t,
		Primitive: p,
	}
}

func writeChar(r rune) string {
	switch r {
	case '\n':
		return "#\\newline"
	case ' ':
		return "#\\space"
	case '\t':
		return "#\\tab"
	// add others if you like
	default:
		// If you want #\a for printable runes:
		return "#\\" + string(r)
	}
}
