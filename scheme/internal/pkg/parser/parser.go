package parser

import (
	"bytes"
	"context"
	"fmt"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
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

func (p *Parser) SetPrompt(prompt string) {
	p.prompt = prompt
}

func (p *Parser) Eval(rt *Runtime) (values.Interface, error) {
	val, err := EvalSExpression(p, rt)
	p.exprnNo++
	_, err = displayImpl(val, rt)
	if err != nil {
		_, _ = fmt.Fprintf(rt.Err, "Error %v\n", err)
	}
	return val, err
}

func (p *Parser) Repl(rtOpts ...OptionRuntime) {

	rt := NewRuntime(rtOpts...)
	p.doPrompt(rt)
	select {
	case <-p.ctx.Done():
		return
	default:
		for tok := p.tokSrc.NextToken(); ; tok = p.tokSrc.NextToken() {
			switch tok.Type {
			case lexer.TokenEOF:
				if p.verbose > Quiet {
					_, _ = fmt.Fprintln(rt.Out, "Bye")
				}
				return
			case lexer.TokenError:
				_, _ = fmt.Fprintf(rt.Err, "Error %#v", tok)
			case lexer.TokenLParen:
				//start new S - Expression
				val, err := EvalSExpression(p, rt)
				if err != nil {
					_, _ = fmt.Fprintf(rt.Err, "Error %v\n", err)
				}
				_, err = displayImpl(val, rt)
				if err != nil {
					_, _ = fmt.Fprintf(rt.Err, "Error %v\n", err)
				}
				p.doPrompt(rt)
			}
			p.exprnNo++
		}
	}
}

func (p *Parser) doPrompt(rt *Runtime) {
	if p.showExpressionCount && p.prompt != "" {
		_, _ = fmt.Fprintf(rt.Out, "%d:%s", p.exprnNo, p.prompt)
	} else if p.prompt != "" {
		_, _ = fmt.Fprintf(rt.Out, "%s ", p.prompt)
	}
}

func EvalString(ctx context.Context, str string, rt *Runtime) (values.Interface, error) {
	p := New(ctx, lexer.New(bytes.NewBufferString(str)))
	val, err := EvalSExpression(p, rt)
	p.exprnNo++
	_, err = displayImpl(val, rt)
	if err != nil {
		_, _ = fmt.Fprintf(rt.Err, "Error %v\n", err)
	}
	return val, err
}

func ReadDatum(p *Parser, rt *Runtime) (values.Interface, error) {
	select {
	case <-p.ctx.Done():
		return values.NewVoidType(), p.ctx.Err()
	default:
		tok := p.tokSrc.NextToken()

		var atoms = values.NewNil()
		for ; tok.Type != lexer.TokenEOF; tok = p.tokSrc.NextToken() {
			if p.verbose >= Debug {
				_, _ = fmt.Fprintf(rt.Err, "Token Runtime: %v Token Literal: %v\n", tok.Type, tok.Literal)
			}
			switch tok.Type {
			case lexer.TokenQuot:
				quotedExpr, err := ReadDatum(p, rt)
				if err != nil {
					return values.NewVoidType(), err
				}
				//quotPair := values.Cons(values.NewIdentifier("quot"), values.Cons(quotedExpr, values.NewNil()))
				return values.Cons(values.NewQuotType(), quotedExpr), nil
			case lexer.TokenEOF:
				return values.NewVoidType(), nil
			case lexer.TokenError:
				return values.NewVoidType(), ErrInvalidToken
			case lexer.TokenLParen:
				nestedExpr, err := ReadDatum(p, rt)
				if err != nil {
					return values.NewVoidType(), err
				}
				if atoms.Type() == types.Nil {
					atoms = nestedExpr
					continue
				}
				atoms = values.Cons(nestedExpr, atoms)
			case
				lexer.TokenIdent,
				lexer.TokenInt,
				lexer.TokenString,
				lexer.TokenBoolean,
				lexer.TokenRelationalOperator,
				lexer.TokenArithmeticOperator:
				atoms = values.Cons(values.FromToken(tok), atoms)
			case lexer.TokenRParen:
				if atoms.Type() == types.Nil {
					return values.NewNil(), nil
				}
				return values.Reverse(atoms), nil
			default:
				return values.NewVoidType(), ErrInvalidToken
			}
		}
		if values.Cdr(atoms).Type() != types.Nil {
			return atoms, nil
		}
		return values.Car(atoms), nil
	}
}

func EvalSExpression(p *Parser, rt *Runtime) (values.Interface, error) {
	val, err := ReadDatum(p, rt)
	if err != nil {
		return values.NewVoidType(), err
	}
	return evalSexpression(val, rt)
}
