package lexer

import (
	"fmt"
	"text/scanner"
)

type LexError struct {
	Position scanner.Position
	Message  string
}

func (e LexError) Error() string {
	return fmt.Sprintf("error %s at position: %v", e.Message, e.Position)
}
