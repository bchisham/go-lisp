package parser

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

func TestEvalSExpression(t *testing.T) {
	type args struct {
		p *Parser
	}
	tests := []struct {
		name    string
		args    args
		want    types.Value
		wantErr bool
	}{
		{
			name: "hello world",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(format t \"hello world\")"))),
			},
			want: newVoidType(),
		},
		{
			name: "display",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(display \"hello world\" (newline))"))),
			},
			want: newVoidType(),
		},
		{
			name: "quot",
			args: args{
				p: New(context.Background(), lexer.New(bytes.NewBufferString("(quot (1 2 3 4))"))),
			},
			want: newList(
				newInt(1),
				newInt(2),
				newInt(3),
				newInt(4)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvalSExpression(tt.args.p, defaultEnvironment())
			if (err != nil) != tt.wantErr {
				t.Errorf("EvalSExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvalSExpression() got = %v, want %v", got, tt.want)
			}
		})
	}
}
