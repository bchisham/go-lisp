package values

import (
	"strings"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type Pair interface {
	Interface
	Car() Interface
	Cdr() Interface
}

type pairVal struct {
	truthyValue
	car Interface
	cdr Interface
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
		car: car,
		cdr: cdr,
	}
}

func Car(p Interface) Interface {
	pair, ok := p.(pairVal)
	if !ok {
		panic("car called on non-pair")
	}
	return pair.Car()
}

func Cdr(p Interface) Interface {
	pair, ok := p.(pairVal)
	if !ok {
		panic("cdr called on non-pair")
	}
	return pair.Cdr()
}

func (pr pairVal) Car() Interface {
	return pr.car
}

func (pr pairVal) Cdr() Interface {
	return pr.cdr
}

func Reverse(input Interface) (output Interface) {
	output = NewNil()
	current := input
	for {
		pair, ok := current.(pairVal)
		if !ok {
			break
		}
		if pair.Cdr().Type() == types.Nil {
			return Cons(pair.Car(), output)
		}
		output = Cons(pair.Car(), output)
		current = pair.Cdr()
	}
	return output
}

func (pr pairVal) Equal(p Interface) bool {
	otherPair, ok := p.(pairVal)
	if !ok {
		return false
	}
	if !pr.Car().Equal(otherPair.Car()) {
		return false
	}
	if !pr.Cdr().Equal(otherPair.Cdr()) {
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

func (pr pairVal) DisplayString() string {

	if pr.Cdr().Type() == types.Nil {
		return pr.Car().DisplayString()
	}

	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(pr.Car().DisplayString())

	cdr := pr.Cdr()
	for {
		if _, ok := cdr.(Nil); ok {
			break
		}
		if pair, ok := cdr.(pairVal); ok {
			sb.WriteString(" ")
			sb.WriteString(pair.Car().DisplayString())
			cdr = pair.Cdr()
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

	if pr.Cdr().Type() == types.Nil {
		return pr.Car().WriteString()
	}

	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(pr.Car().WriteString())

	cdr := pr.Cdr()
	for {
		if _, ok := cdr.(Nil); ok {
			break
		}
		if pair, ok := cdr.(pairVal); ok {
			sb.WriteString(" ")
			sb.WriteString(pair.Car().WriteString())
			cdr = pair.Cdr()
		} else {
			sb.WriteString(" . ")
			sb.WriteString(cdr.WriteString())
			break
		}
	}
	sb.WriteString(")")
	return sb.String()
}
