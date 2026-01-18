package parser

import (
	"bytes"
	"lisp/internal/pkg/lexer"
	"lisp/internal/pkg/parser/types"
	"reflect"
	"testing"
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
				p: New(lexer.New(bytes.NewBufferString("(format t \"hello world\")"))),
			},
			want: newVoidType(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvalSExpression(tt.args.p)
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
