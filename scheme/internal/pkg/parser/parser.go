package parser

import (
	"context"
	"fmt"
	"os"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

const (
	Quiet VerboseLevel = iota
	Error VerboseLevel = iota
	Warn  VerboseLevel = iota
	Info  VerboseLevel = iota
	Debug VerboseLevel = iota
)

type VerboseLevel int

type config struct {
	prompt              string
	verbose             VerboseLevel
	showExpressionCount bool
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

func WithShowExpressionCount(showExpressionCount bool) Option {
	return func(c *config) {
		c.showExpressionCount = showExpressionCount
	}
}

type Parser struct {
	config
	ctx     context.Context
	tokSrc  *lexer.Scanner
	exprnNo int
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
	p.doPrompt()
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
					fmt.Printf("Error %v\n", err)
				}
				_, _ = displayImpl(list.New(val), env)
				fmt.Printf("%s", p.prompt)
			}
			p.exprnNo++
		}
	}
}

func (p *Parser) doPrompt() {
	if p.showExpressionCount && p.prompt != "" {
		fmt.Printf("%d:%s", p.exprnNo, p.prompt)
	} else if p.prompt != "" {
		fmt.Printf("%s ", p.prompt)
	}
}

func EvalSExpression(p *Parser, env values.Environment) (values.Value, error) {
	tok := p.tokSrc.NextToken()

	var atoms []values.Value
	for ; tok.Type != lexer.TokenEOF; tok = p.tokSrc.NextToken() {
		if p.verbose >= Debug {
			_, _ = fmt.Fprintf(os.Stderr, "Token Type: %v Token Literal: %v\n", tok.Type, tok.Literal)
		}
		switch tok.Type {
		case lexer.TokenEOF:
			return values.NewVoidType(), nil
		case lexer.TokenError:
			return values.NewVoidType(), ErrInvalidToken
		case lexer.TokenLParen:
			nestedExpr, err := EvalSExpression(p, env)
			if err != nil {
				return values.NewVoidType(), err
			}
			atoms = append(atoms, nestedExpr)
		case lexer.TokenQuot,
			lexer.TokenIdent,
			lexer.TokenInt,
			lexer.TokenString,
			lexer.TokenBoolean,
			lexer.TokenRelationalOperator,
			lexer.TokenArithmeticOperator:
			atoms = append(atoms, values.FromToken(tok))
		case lexer.TokenRParen:
			goto eval
		}
	}
eval:
	return evalSexpression(atoms, env)
}
