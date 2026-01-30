package lexer

import (
	"fmt"
	"slices"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer/types"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/list"

	"io"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/boolean"
)

type TokenType string

const (
	TokenEOF                TokenType = "EOF"
	TokenError              TokenType = "error"
	TokenNumber             TokenType = "number"
	TokenInt                TokenType = "int"
	TokenString             TokenType = "string"
	TokenBoolean            TokenType = "boolean"
	TokenRune               TokenType = "rune"
	TokenPair               TokenType = "pair"
	TokenSymbol             TokenType = "symbol"
	TokenIdent              TokenType = "ident"
	TokenColonIdent         TokenType = "colon_ident"
	TokenLParen             TokenType = "("
	TokenRParen             TokenType = ")"
	TokenLBracket           TokenType = "["
	TokenRBracket           TokenType = "]"
	TokenLBrace             TokenType = "{"
	TokenRBrace             TokenType = "}"
	TokenSemiColon          TokenType = ";"
	TokenQuot               TokenType = "'"
	TokenLineComment        TokenType = "line_comment"
	TokenRelationalOperator TokenType = "relationalOperator"
	TokenArithmeticOperator TokenType = "arithmeticOperator"
)

type Token struct {
	Type    TokenType
	Literal string
	Int     int64
	Float   float64
	Bool    bool
	Text    string
	Ident   string
	Error   LexError
}

func (t Token) String() string {
	return t.Literal
}

type Scanner struct {
	scan scanner.Scanner
}

func New(r io.Reader) *Scanner {
	s := &Scanner{
		scanner.Scanner{},
	}
	s.scan.Init(r)
	return s
}

// NextToken extract the next token from the input stream
func (s *Scanner) NextToken() (tok Token) {

	for ch := s.scan.Peek(); ch != scanner.EOF; ch = s.scan.Peek() {
		if unicode.IsDigit(ch) {
			return s.consumeNumber()
		}
		if unicode.IsLetter(ch) {
			return s.consumeIdentifier()
		}
		if unicode.IsSpace(ch) {
			s.scan.Next()
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
		case '\'':
			return s.consumeQuot()
		case '<', '>', '=':
			return s.consumeRelationalOperator()
		case '+', '-', '*', '/', '%':
			return s.consumeArithmeticOperatorOrNumber()
		case '#':
			return s.consumeLiteral()
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
		Type:    TokenInt,
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
		quotChar = s.scan.Next()
	)

	if quotChar == scanner.String {
		quotChar = '"'
	}
	//sb.WriteRune(quotChar)

	for ch := s.scan.Peek(); ch != scanner.EOF && ch != scanner.String; ch = s.scan.Peek() {
		if ch == quotChar {
			s.scan.Next()
			break
		}
		content.WriteRune(s.scan.Next())
		sb.WriteRune(ch)
	}
	//sb.WriteRune(quotChar)

	return Token{
		Type:    TokenString,
		Literal: sb.String(),
		Text:    content.String(),
	}

}

func (s *Scanner) consumeLParen() (tok Token) {
	_ = s.scan.Next()
	return Token{
		Type:    TokenLParen,
		Literal: "(",
	}
}

func (s *Scanner) consumeRParen() (tok Token) {
	_ = s.scan.Next()
	return Token{
		Type:    TokenRParen,
		Literal: ")",
	}
}

var (
	relationalOperRunes      = list.New('<', '>', '=')
	relationalOperStartsWith = boolean.AnyFunc(onlyRunes(relationalOperRunes))
	relationalOperEndToken   = boolean.NotFunc(relationalOperStartsWith)
)

func (s *Scanner) consumeColonIdent() Token {
	s.scan.Next()
	txt := s.collectRunes(startIdentifierFunc, continueIdentifierFunc)
	return Token{
		Type:    TokenColonIdent,
		Literal: ":" + txt,
		Ident:   txt,
	}
}

// line-comments start with ";;" followed by any text to end of the line
func (s *Scanner) consumeSemiColonOrLineComment() Token {
	s.scan.Next()
	if s.scan.Peek() != ';' {
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

func (s *Scanner) consumeQuot() Token {
	_ = s.scan.Next()
	return Token{
		Type:    TokenQuot,
		Literal: string(TokenQuot),
	}
}

func (s *Scanner) consumeRelationalOperator() Token {
	start := s.scan.Next()
	if slices.Contains(list.New('<', '>'), start) && '=' == s.scan.Peek() {
		next := s.scan.Next()
		return Token{
			Type:    TokenRelationalOperator,
			Literal: string(start) + string(next),
		}
	}
	return Token{
		Type:    TokenRelationalOperator,
		Literal: string(start),
	}
}

func (s *Scanner) consumeArithmeticOperatorOrNumber() Token {
	oper := s.scan.Next()
	if oper == '-' && unicode.IsDigit(s.scan.Peek()) {
		//this is a number not ([operator]  .... )
		number := s.consumeNumber()
		number.Literal = fmt.Sprintf("%c%s", oper, number)
		number.Int = -number.Int
		number.Float = -number.Float
		return number
	}
	return Token{
		Type:    TokenArithmeticOperator,
		Literal: string(oper),
	}
}

func (s *Scanner) consumeLiteral() Token {
	_ = s.scan.Next() //#
	if slices.Contains(list.New('t', 'f'), s.scan.Peek()) {
		val := s.scan.Next()
		return Token{
			Type:    TokenBoolean,
			Literal: "#" + string(val),
			Bool:    val == 't',
		}
	} else if s.scan.Peek() == '\\' {
		_ = s.scan.Next() //consume the '\'
		charLiteral := s.scan.Next()

		return Token{
			Type:    TokenRune,
			Literal: fmt.Sprintf("#\\%c", charLiteral),
			Text:    string(charLiteral),
			Int:     int64(charLiteral),
		}
	}
	return Token{
		Type:    TokenError,
		Literal: "#" + string(s.scan.Peek()),
		Error: LexError{
			Message: "invalid boolean literal",
		},
	}
}

func (s *Scanner) collectRunes(startsWith types.RuneClassifier, exitCondition types.RuneClassifier) string {
	var sb strings.Builder
	if !startsWith(s.scan.Peek()) {
		return ""
	}
	sb.WriteRune(s.scan.Next())

	for ch := s.scan.Peek(); ch != scanner.EOF; ch = s.scan.Peek() {
		if exitCondition(ch) {
			goto tokenComplete
		}
		sb.WriteRune(s.scan.Next())

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
