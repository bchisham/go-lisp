package parser

import "errors"

var (
	ErrOperatorIsNotAProcedure = errors.New("operator is not a procedure")
	ErrUnexpectedToken         = errors.New("unexpected token")
	ErrUndefinedIdent          = errors.New("undefined identifier")
	ErrInvalidToken            = errors.New("invalid token")
	ErrEof                     = errors.New("eof")
	ErrWrongNumberOfArguments  = errors.New("wrong number of arguments")
	ErrNumberExpected          = errors.New("number expected")
)
