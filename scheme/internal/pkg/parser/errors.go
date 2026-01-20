package parser

import (
	"errors"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

var (
	ErrInvalidFormat           = errors.New("invalid format")
	ErrBadArgument             = errors.New("bad argument")
	ErrOperatorIsNotAProcedure = errors.New("operator is not a procedure")
	ErrNotAPrimitive           = values.ErrNotAPrimitive
	ErrUnexpectedToken         = errors.New("unexpected token")
	ErrUndefinedIdent          = errors.New("undefined identifier")
	ErrInvalidToken            = errors.New("invalid token")
	ErrEof                     = errors.New("eof")
	ErrWrongNumberOfArguments  = errors.New("wrong number of arguments")
	ErrNumberExpected          = errors.New("number expected")
)

func ErrIo(err error) error {
	return ErrType{
		err:     err,
		message: "I/O error",
	}
}

type ErrType struct {
	message string
	err     error
}

func (e ErrType) Error() string {
	return e.message
}

func (e ErrType) Unwrap() error {
	return e.err
}
