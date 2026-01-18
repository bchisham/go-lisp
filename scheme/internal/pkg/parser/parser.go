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
	fmt.Printf("%s", p.prompt)
	env := defaultEnvironment()
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
				val, err := EvalSExpression(p, env)
				if err != nil {
					return
				}
				_, _ = displayImpl(list.New(val), env)
				fmt.Printf("%s", p.prompt)
			}
		}
	}
}

func EvalSExpression(p *Parser, env types.Environment) (types.Value, error) {
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
			nestedExpr, err := EvalSExpression(p, env)
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
	return evalSexpression(atoms, env)
}
