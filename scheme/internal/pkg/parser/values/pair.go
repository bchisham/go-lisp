package values

import (
	"strings"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type pairVal struct {
	Car Interface
	Cdr Interface
}

func Cons(car, cdr Interface) Interface {
	if car == nil {
		panic("car cannot be nil")
	}
	if cdr == nil {
		panic("cdr cannot be nil")
	}
	if car.Type() == types.Nil && cdr.Type() == types.Nil {
		return car
	}
	if car.Type() == types.Nil && cdr.Type() != types.Nil {
		panic("car cannot be nil if cdr is not nil")
	}
	return pairVal{
		Car: car,
		Cdr: cdr,
	}
}

func Car(p Interface) Interface {
	if pv, ok := p.(pairVal); ok {
		return pv.Car
	}
	return NilValue{}
}

func Cdr(p Interface) Interface {
	if pv, ok := p.(pairVal); ok {
		return pv.Cdr
	}
	return NilValue{}
}

func Reverse(input Interface) (output Interface) {
	output = NewNil()
	current := input
	for {
		pair, ok := current.(pairVal)
		if !ok {
			break
		}
		if pair.Cdr.Type() == types.Nil {
			return Cons(pair.Car, output)
		}
		output = Cons(pair.Car, output)
		current = pair.Cdr
	}
	return output
}

func (pr pairVal) Equal(p Interface) bool {
	otherPair, ok := p.(pairVal)
	if !ok {
		return false
	}
	if !pr.Car.Equal(otherPair.Car) {
		return false
	}
	if !pr.Cdr.Equal(otherPair.Cdr) {
		return false
	}
	return true

}

func (pr pairVal) Type() types.Type {

	return types.Pair
}

func (pr pairVal) GetToken() lexer.Token {
	return lexer.Token{
		Type: lexer.TokenPair,
	}
}

func (pr pairVal) AsPrimitive() (Primitive, error) {
	return Primitive{}, ErrNotAPrimitive
}

func (pr pairVal) IsTruthy() bool {
	return true
}

func (pr pairVal) DisplayString() string {

	if pr.Cdr.Type() == types.Nil {
		return pr.Car.DisplayString()
	}

	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(pr.Car.DisplayString())

	cdr := pr.Cdr
	for {
		if _, ok := cdr.(NilValue); ok {
			break
		}
		if pair, ok := cdr.(pairVal); ok {
			sb.WriteString(" ")
			sb.WriteString(pair.Car.DisplayString())
			cdr = pair.Cdr
		} else {
			sb.WriteString(" . ")
			sb.WriteString(cdr.DisplayString())
			break
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func (pr pairVal) WriteString() string {

	if pr.Cdr.Type() == types.Nil {
		return pr.Car.WriteString()
	}

	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(pr.Car.WriteString())

	cdr := pr.Cdr
	for {
		if _, ok := cdr.(NilValue); ok {
			break
		}
		if pair, ok := cdr.(pairVal); ok {
			sb.WriteString(" ")
			sb.WriteString(pair.Car.WriteString())
			cdr = pair.Cdr
		} else {
			sb.WriteString(" . ")
			sb.WriteString(cdr.WriteString())
			break
		}
	}
	sb.WriteString(")")
	return sb.String()
}
