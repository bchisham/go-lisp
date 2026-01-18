package parser

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

var ErrUnexpectedToken = errors.New("unexpected token")
var ErrInvalidToken = errors.New("invalid token")
var ErrUndefinedIdent = errors.New("undefined identifier")
var ErrEof = errors.New("eof")

const (
	Quiet VerboseLevel = iota
	Error VerboseLevel = iota
	Warn  VerboseLevel = iota
	Info  VerboseLevel = iota
	Debug VerboseLevel = iota
)

type VerboseLevel int

type config struct {
	prompt  string
	verbose VerboseLevel
}

type Option func(*config)

func WithVerbose(verbose VerboseLevel) Option {
	return func(c *config) {
		c.verbose = verbose
	}
}

func WithPrompt(prompt string) Option {
	return func(c *config) {
		c.prompt = prompt
	}
}

type Parser struct {
	config
	ctx    context.Context
	tokSrc *lexer.Scanner
}

func New(ctx context.Context, tokSrc *lexer.Scanner, opts ...Option) *Parser {
	cfg := config{
		prompt: "go-scheme> ",
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Parser{
		config: cfg,
		ctx:    ctx,
		tokSrc: tokSrc,
	}
}

func (p *Parser) Repl() {
	fmt.Printf(p.prompt)

	select {
	case <-p.ctx.Done():
		return
	default:
		for tok := p.tokSrc.NextToken(); ; tok = p.tokSrc.NextToken() {
			switch tok.Type {
			case lexer.TokenEOF:
				if p.verbose > Quiet {
					fmt.Println("Bye")
				}
				return
			case lexer.TokenError:
				fmt.Printf("Error %#v", tok)
			case lexer.TokenLParen:
				//start new S - Expression
				val, err := EvalSExpression(p)
				if err != nil {
					return
				}
				printValue(val)
				fmt.Printf(p.prompt)
			}
		}
	}
}

func printValue(val types.Value) {
	switch val.Type {
	case types.String:
		fmt.Printf("%s", val.StringVal)
	case types.Int:
		fmt.Printf("%d", val.IntVal)
	case types.Float:
		fmt.Printf("%f", val.FloatVal)
	case types.Bool:
		fmt.Printf("%t", val.BoolVal)
	}
	fmt.Println()
}

func EvalSExpression(p *Parser) (types.Value, error) {
	tok := p.tokSrc.NextToken()

	var atoms []types.Value
	for ; tok.Type != lexer.TokenEOF; tok = p.tokSrc.NextToken() {
		if p.verbose >= Debug {
			_, _ = fmt.Fprintf(os.Stderr, "Token Type: %v Token Literal: %v\n", tok.Type, tok.Literal)
		}
		switch tok.Type {
		case lexer.TokenEOF:
			return newVoidType(), nil
		case lexer.TokenError:
			return newVoidType(), ErrInvalidToken
		case lexer.TokenLParen:
			nestedExpr, err := EvalSExpression(p)
			if err != nil {
				return newVoidType(), err
			}
			atoms = append(atoms, nestedExpr)
		case lexer.TokenIdent:
			ident := newIdentifier(tok.Literal)
			atoms = append(atoms, ident)
		case lexer.TokenInt:
			atoms = append(atoms, newInt(tok.Int))
		case lexer.TokenString:
			atoms = append(atoms, newString(tok.Text))
		case lexer.TokenRParen:
			goto eval
		}
	}
eval:
	return evalSexpression(atoms)
}

func defaultEnvironment() *types.Environment {
	var builtins = types.NewEnvironment()

	builtins.Define("format", newLambda(builtins, formatImpl))

	return builtins
}

func evalSexpression(l []types.Value) (types.Value, error) {
	head := list.Car(l)
	tail := list.Cdr(l)
	env := defaultEnvironment()
	switch head.Type {
	case types.Lambda:
		return head.LambdaVal.Body(tail)
	case types.Identifier:
		resvVal, ok := env.Lookup(head.NameVal)
		if !ok {
			return newVoidType(), ErrUndefinedIdent
		}
		switch resvVal.Type {
		case types.Lambda:
			return resvVal.LambdaVal.Body(tail)

		default:
			return resvVal, nil
		}
	case types.List:
		//TODO handle head is list
	default:
		return head, nil

	}
	return newVoidType(), nil
}

func newVoidType() types.Value {
	return types.Value{
		Type: types.Void,
	}
}

func newString(s string) types.Value {
	return types.Value{
		Type: types.String,
		Primitive: types.Primitive{
			StringVal: s,
		},
	}
}

func newFloat(f float64) types.Value {
	return types.Value{
		Type: types.Float,
		Primitive: types.Primitive{
			FloatVal: f,
		},
	}
}

func newInt(i int64) types.Value {
	return types.Value{
		Type: types.Int,
		Primitive: types.Primitive{
			IntVal: i,
		},
	}
}

func newBool(b bool) types.Value {
	return types.Value{
		Type: types.Bool,
		Primitive: types.Primitive{
			BoolVal: b,
		},
	}
}
func newChar(c rune) types.Value {
	return types.Value{
		Type: types.Char,
		Primitive: types.Primitive{
			CharVal: string(c),
		},
	}
}

func newLambda(env *types.Environment, expression types.Expression) types.Value {
	return types.Value{
		Type: types.Lambda,
		LambdaVal: types.LambdaExpr{
			Env:  env,
			Body: expression,
		},
	}
}

func newExpression(env types.Environment, body []types.Value) types.Expression {
	return func(args []types.Value) (types.Value, error) {
		return newVoidType(), nil
	}
}

func newIdentifier(name string) types.Value {
	return types.Value{
		Type: types.Identifier,
		Primitive: types.Primitive{
			NameVal: name,
		},
	}
}
