package parser

import (
	"bytes"
	"context"
	"testing"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

func TestEvalSExpression(t *testing.T) {
	type args struct {
		p *Parser
	}
	tests := []struct {
		name    string
		args    args
		want    values.Value
		wantErr bool
	}{
		{
			name: "hello world",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(format t \"hello world\")"))),
			},
			want: values.NewVoidType(),
		},
		{
			name: "display",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(display \"hello world\" (newline))"))),
			},
			want: values.NewVoidType(),
		},
		{
			name: "quot - list",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("'(1 2 3 4)"))),
			},
			want: values.NewList(
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4)),
		},
		{
			name: "quot - literal",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("'1"))),
			},
			want: values.NewInt(1),
		},
		{
			name: "less-than: expect true two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(< 1 2)"))),
			},
			want: values.NewBool(true),
		},
		{
			name: "less-than: expect true three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(< 1 2 3)"))),
			},
			want: values.NewBool(true),
		},
		{
			name: "less-than: expect false two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(< 4 2)"))),
			},
			want: values.NewBool(false),
		},
		{
			name: "less-than: expect false three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(< 4 2 3)"))),
			},
			want: values.NewBool(false),
		},
		{
			name: "greater-than: expect true two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(> 2 1)"))),
			},
			want: values.NewBool(true),
		},
		{
			name: "greater-than: expect true three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(> 2 1 0)"))),
			},
			want: values.NewBool(true),
		},
		{
			name: "greater-than: expect false two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(> 2 4)"))),
			},
			want: values.NewBool(false),
		},
		{
			name: "greater-than: expect false three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(> 4 2 3)"))),
			},
			want: values.NewBool(false),
		},
		{
			name: "equal: expect true two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(= 2 2)"))),
			},
			want: values.NewBool(true),
		},
		{
			name: "equal: expect true three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(= 2 2 2)"))),
			},
			want: values.NewBool(true),
		},
		{
			name: "equal: expect false three operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(= 2 2 3)"))),
			},
			want: values.NewBool(false),
		},
		{
			name: "equal: expect false two operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(= 2 3)"))),
			},
			want: values.NewBool(false),
		},
		{
			name: "addition: no operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(+)"))),
			},
			want: values.NewInt(0),
		},
		{
			name: "addition: 1 operand",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(+ 1)"))),
			},
			want: values.NewInt(1),
		},
		{
			name: "addition: 2 operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(+ 1 2)"))),
			},
			want: values.NewInt(3),
		},
		{
			name: "addition: 3 operands",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(+ 1 2 3)"))),
			},
			want: values.NewInt(6),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvalSExpression(tt.args.p, defaultEnvironment())
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
