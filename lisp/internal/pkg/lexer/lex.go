package lexer

import (
	"lisp/internal/pkg/lexer/types"

	"io"
	"lisp/internal/pkg/boolean"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
)

const (
	KwAnd    = "and"
	KwOr     = "or"
	KwXor    = "xor"
	KwNot    = "not"
	KwCons   = "cons"
	KwCdr    = "cdr"
	KwEq     = "eq"
	KwLt     = "lt"
	KwGt     = "gt"
	KwLe     = "le"
	KwGtEq   = "gt_eq"
	KwLtEq   = "lt_eq"
	KwLambda = "lambda"
	KwIf     = "if"
)

type TokenType string

const (
	TokenEOF         TokenType = "EOF"
	TokenError       TokenType = "error"
	TokenNumber      TokenType = "number"
	TokenString      TokenType = "string"
	TokenSymbol      TokenType = "symbol"
	TokenIdent       TokenType = "ident"
	TokenColonIdent  TokenType = "colon_ident"
	TokenLParen      TokenType = "("
	TokenRParen      TokenType = ")"
	TokenLBracket    TokenType = "["
	TokenRBracket    TokenType = "]"
	TokenLBrace      TokenType = "{"
	TokenRBrace      TokenType = "}"
	TokenSemiColon   TokenType = ";"
	TokenLineComment TokenType = "line_comment"
)

type Token struct {
	Type    TokenType
	Literal string
	Int     int64
	Float   float64
	Text    string
	Ident   string
	Error   LexError
}

func (t Token) String() string {
	return t.Literal
}

type Scanner struct {
	scanner.Scanner
}

func New(r io.Reader) *Scanner {
	s := &Scanner{
		scanner.Scanner{},
	}
	s.Init(r)
	return s
}

func (s *Scanner) NextToken() (tok Token) {

	for ch := s.Peek(); ch != scanner.EOF; ch = s.Peek() {
		if unicode.IsDigit(ch) || ch == '-' {
			return s.consumeNumber()
		}
		if unicode.IsLetter(ch) {
			return s.consumeIdentifier()
		}
		switch ch {
		case '[':
			return Token{
				Type:    TokenLBracket,
				Literal: "[",
			}
		case ']':
			return Token{
				Type:    TokenRBracket,
				Literal: "]",
			}
		case '(':
			return s.consumeLParen()
		case ')':
			return s.consumeRParen()
		case ':':
			return s.consumeColonIdent()
		case '"':
			return s.consumeString()
		}

	}

	return Token{Type: TokenEOF}
}

var (
	startNumberFunc    = boolean.AnyFunc(onlyRunes([]rune{'-'}), isNumberChar)
	continueNumberFunc = boolean.NotFunc(boolean.AnyFunc(unicode.IsDigit, onlyRunes([]rune{'.'})))
)

func (s *Scanner) consumeNumber() (_ Token) {
	itxt := s.collectRunes(startNumberFunc, continueNumberFunc)
	itxt = strings.TrimSpace(itxt)
	intval, err := strconv.ParseInt(itxt, 10, 64)
	if err != nil {
		return Token{
			Type:    TokenError,
			Literal: itxt,
		}
	}

	return Token{
		Type:    TokenNumber,
		Literal: itxt,
		Int:     intval,
	}
}

var (
	startIdentifierFunc    = unicode.IsLetter
	continueIdentifierFunc = boolean.NotFunc(isIdentifierChar)
)

func (s *Scanner) consumeIdentifier() (_ Token) {

	sb := strings.Builder{}
	sb.WriteString(s.collectRunes(startIdentifierFunc, continueIdentifierFunc))

	return Token{
		Type:    TokenIdent,
		Literal: sb.String(),
		Ident:   sb.String(),
	}
}

func (s *Scanner) consumeString() (_ Token) {
	var (
		sb       strings.Builder
		content  strings.Builder
		quotChar = s.Next()
	)

	if quotChar == scanner.String {
		quotChar = '"'
	}
	sb.WriteRune(quotChar)

	for ch := s.Peek(); ch != scanner.EOF && ch != scanner.String; ch = s.Peek() {
		if ch == quotChar {
			break
		}
		content.WriteRune(s.Next())
		sb.WriteRune(ch)
	}
	sb.WriteRune(quotChar)

	return Token{
		Type:    TokenString,
		Literal: sb.String(),
		Text:    content.String(),
	}

}

func (s *Scanner) consumeLParen() (tok Token) {
	_ = s.Next()
	return Token{
		Type:    TokenLParen,
		Literal: "(",
	}
}

func (s *Scanner) consumeRParen() (tok Token) {
	_ = s.Next()
	return Token{
		Type:    TokenRParen,
		Literal: ")",
	}
}

func (s *Scanner) consumeColonIdent() Token {
	s.Next()
	txt := s.collectRunes(startIdentifierFunc, continueIdentifierFunc)
	return Token{
		Type:    TokenColonIdent,
		Literal: ":" + txt,
		Ident:   txt,
	}
}

// line-comments start with ";;" followed by any text to end of the line
func (s *Scanner) consumeSemiColonOrLineComment() Token {
	s.Next()
	if s.Peek() != ';' {
		return Token{
			Type:    TokenSemiColon,
			Literal: ";",
		}
	}
	_ = s.collectRunes(unicode.IsGraphic, anyBut([]rune{'\n'}))
	return Token{
		Type: TokenLineComment,
	}
}

func (s *Scanner) collectRunes(startsWith types.RuneClassifier, exitCondition types.RuneClassifier) string {
	var sb strings.Builder
	if !startsWith(s.Peek()) {
		return ""
	}
	sb.WriteRune(s.Next())

	for ch := s.Peek(); ch != scanner.EOF; ch = s.Peek() {
		if exitCondition(ch) {
			goto tokenComplete
		}
		sb.WriteRune(s.Next())

	}
tokenComplete:
	return sb.String()
}

func isTerm(termChar rune) types.RuneClassifier {
	return func(r rune) bool {
		return termChar == r
	}
}

func anyStart(rune) bool { return true }

func anyBut(notAllowed []rune) types.RuneClassifier {
	return func(r rune) bool {
		for _, a := range notAllowed {
			if a == r {
				return true
			}
		}
		return false
	}
}
func onlyRunes(allowed []rune) types.RuneClassifier {
	return func(r rune) bool {
		for _, a := range allowed {
			if a == r {
				return true
			}
		}
		return false
	}
}

func isIdentifierChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'
}

func isNumberChar(c rune) bool {
	return unicode.IsDigit(c) || c == '.'
}
