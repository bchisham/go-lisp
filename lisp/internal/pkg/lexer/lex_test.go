package lexer

import (
	"bytes"
	"reflect"
	"testing"
)

func TestScanner_NextToken(t *testing.T) {
	type fields struct {
		Scanner *Scanner
	}
	tests := []struct {
		name    string
		fields  fields
		wantTok Token
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				Scanner: New(bytes.NewBuffer(nil)),
			},
			wantErr: false,
			wantTok: Token{
				Type: TokenEOF,
			},
		},
		{
			name: "number",
			fields: fields{
				Scanner: New(bytes.NewBuffer([]byte("123456"))),
			},
			wantErr: false,
			wantTok: Token{
				Type:    TokenNumber,
				Literal: "123456",
				Int:     123456,
			},
		},
		{
			name: "left-paren",
			fields: fields{
				Scanner: New(bytes.NewBuffer([]byte("("))),
			},
			wantErr: false,
			wantTok: Token{
				Type:    TokenLParen,
				Literal: "(",
			},
		},
		{
			name: "right-paren",
			fields: fields{
				Scanner: New(bytes.NewBuffer([]byte(")"))),
			},
			wantErr: false,
			wantTok: Token{
				Type:    TokenRParen,
				Literal: ")",
			},
		},
		{
			name: "left-bracket",
			fields: fields{
				Scanner: New(bytes.NewBuffer([]byte("["))),
			},
			wantErr: false,
			wantTok: Token{
				Type:    TokenLBracket,
				Literal: "[",
			},
		},
		{
			name: "right-bracket",
			fields: fields{
				Scanner: New(bytes.NewBuffer([]byte("]"))),
			},
			wantErr: false,
			wantTok: Token{
				Type:    TokenRBracket,
				Literal: "]",
			},
		},
		{
			name: "identifier",
			fields: fields{
				Scanner: New(bytes.NewBuffer([]byte("foo"))),
			},
			wantErr: false,
			wantTok: Token{
				Type:    TokenIdent,
				Literal: "foo",
				Ident:   "foo",
			},
		},
		{
			name: "colon identifier",
			fields: fields{
				Scanner: New(bytes.NewBuffer([]byte(":foo"))),
			},
			wantErr: false,
			wantTok: Token{
				Type:    TokenColonIdent,
				Literal: ":foo",
				Ident:   "foo",
			},
		},
		{
			name: "string",
			fields: fields{
				Scanner: New(bytes.NewBuffer([]byte("\"foo\""))),
			},
			wantErr: false,
			wantTok: Token{
				Type:    TokenString,
				Literal: "\"foo\"",
				Text:    "foo",
			},
		},
		{
			name: "extract only one token",
			fields: fields{
				Scanner: New(bytes.NewBuffer([]byte("(lambda"))),
			},
			wantErr: false,
			wantTok: Token{
				Type:    TokenLParen,
				Literal: "(",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotTok := tt.fields.Scanner.NextToken()

			if !reflect.DeepEqual(gotTok, tt.wantTok) {
				t.Errorf("NextToken() gotTok = %v, want %v", gotTok, tt.wantTok)
			}
		})
	}
}
