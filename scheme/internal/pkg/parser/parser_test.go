package parser

import (
	"bytes"
	"context"
	"testing"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/builtins"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

func TestEvalSExpression(t *testing.T) {
	type args struct {
		p  *Parser
		rt *builtins.Runtime
	}
	tests := []struct {
		name    string
		args    args
		want    values.Type
		wantErr bool
	}{
		{
			name: "hello world",
			args: args{
				p: New(
					context.Background(),
					lexer.New(bytes.NewBufferString("(format t \"hello world\")"))),
				rt: builtins.NewRuntime(
					builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewVoidType(),
		},
		{
			name: "display",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(display \"hello world\" (newline))"))),
				rt: builtins.NewRuntime(
					builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewVoidType(),
		},
		{
			name: "quot - list",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("'(1 2 3 4)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewQuotValue(values.Cons(values.NewInt(1),
				values.Cons(values.NewInt(2),
					values.Cons(values.NewInt(3),
						values.Cons(values.NewInt(4), values.NewNil()))))),
		},
		{
			name: "quot - literal",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("'1"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewQuotValue(values.NewInt(1)),
		},
		{
			name: "less-than: expect true two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(< 1 2)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(true),
		},
		{
			name: "less-than: expect true three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(< 1 2 3)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(true),
		},
		{
			name: "less-than: expect false two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(< 4 2)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(false),
		},
		{
			name: "less-than: expect false three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(< 4 2 3)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(false),
		},
		{
			name: "greater-than: expect true two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(> 2 1)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(true),
		},
		{
			name: "greater-than: expect true three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(> 2 1 0)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(true),
		},
		{
			name: "greater-than: expect false two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(> 2 4)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(false),
		},
		{
			name: "greater-than: expect false three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(> 4 2 3)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(false),
		},
		{
			name: "equal: expect true two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(= 2 2)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(true),
		},
		{
			name: "equal: expect true three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(= 2 2 2)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(true),
		},
		{
			name: "equal: expect false three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(= 2 2 3)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(false),
		},
		{
			name: "equal: expect false two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(= 2 3)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewBool(false),
		},
		{
			name: "addition: no operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(+)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(0),
		},
		{
			name: "addition: 1 operand",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(+ 1)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(1),
		},
		{
			name: "addition: 2 operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(+ 1 2)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(3),
		},
		{
			name: "addition: 3 operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(+ 1 2 3)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(6),
		},
		{
			name: "product no operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(*)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(1),
		},
		{
			name: "product 1 operand",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(* 3)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(3),
		},
		{

			name: "product 2 operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(* 3 4)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(12),
		},
		{
			name: "product 3 operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(* 3 4 5)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(60),
		},
		{
			name: "nested expression",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(+ 1 (* 2 3) )"))), //1 + 6 = 7
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(7),
		},
		{
			name: "multiple nested expressions",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(+ 1 (* 2 3) (+ 4 5) )"))), //1 + 6 + 9 = 16
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(16),
		},
		{
			name: "define global variable",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(define x 42)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewVoidType(),
		},
		{
			name: "let binding",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(let ((x 10) (y 20)) (+ x y))"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(30),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvalSExpression(tt.args.p, tt.args.rt)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvalSExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("EvalSExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadDatum(t *testing.T) {
	type args struct {
		p  *Parser
		rt *builtins.Runtime
	}
	tests := []struct {
		name    string
		args    args
		want    values.Type
		wantErr bool
	}{
		{
			name: "read int",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("42"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewInt(42),
		},
		{
			name: "read string",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("\"hello world\""))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewString("hello world"),
		},
		{
			name: "read list",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(1 2 3 4)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.Cons(values.NewInt(1),
				values.Cons(values.NewInt(2),
					values.Cons(values.NewInt(3),
						values.Cons(values.NewInt(4), values.NewNil())))),
		},
		{
			name: "read nested list",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("((1 2) (3 4))"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.Cons(
				values.Cons(values.NewInt(1), values.Cons(values.NewInt(2), values.NewNil())),
				values.Cons(
					values.Cons(values.NewInt(3), values.Cons(values.NewInt(4), values.NewNil())),
					values.NewNil(),
				),
			),
		},
		{
			name: "read quoted list",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("'(1 2 3)"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewQuotValue(values.Cons(values.NewInt(1),
				values.Cons(values.NewInt(2),
					values.Cons(values.NewInt(3), values.NewNil())))),
		},
		{
			name: "read quoted nested list",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("'((1 2) (3 4))"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewQuotValue(values.Cons(
				values.Cons(values.NewInt(1), values.Cons(values.NewInt(2), values.NewNil())),
				values.Cons(
					values.Cons(values.NewInt(3), values.Cons(values.NewInt(4), values.NewNil())),
					values.NewNil(),
				),
			)),
		},
		{
			name: "read empty list",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("()"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.NewNil(),
		},
		{
			name: "read let expression",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(let ((x 10) (y 20)) (+ x y))"))),
				rt: builtins.NewRuntime(builtins.WithOut(bytes.NewBuffer(nil)),
					builtins.WithEvaluatorCallback(evalSexpression)),
			},
			want: values.Cons(
				values.NewIdentifier("let"),
				values.Cons(
					values.Cons(
						values.Cons(values.NewIdentifier("x"), values.Cons(values.NewInt(10), values.NewNil())),
						values.Cons(
							values.Cons(values.NewIdentifier("y"), values.Cons(values.NewInt(20), values.NewNil())),
							values.NewNil(),
						),
					),
					values.Cons(
						values.Cons(values.NewArithmeticOperator("+"), values.Cons(values.NewIdentifier("x"), values.Cons(values.NewIdentifier("y"), values.NewNil()))),
						values.NewNil(),
					),
				),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDatum(tt.args.p, tt.args.rt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDatum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want.Equal(got) {
				t.Errorf("ReadDatum() got = %v, want %v", got, tt.want)
			}
		})
	}
}
