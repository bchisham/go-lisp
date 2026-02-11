package values

import (
	"strings"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

type Pair interface {
	Type
	Car() Type
	Cdr() Type
}

type pairVal struct {
	truthyValue
	car Type
	cdr Type
}

func Cons(car, cdr Type) Type {
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

func Car(p Type) Type {
	pair, ok := p.(pairVal)
	if !ok {
		return NewNil()
	}
	return pair.Car()
}

func Cdr(p Type) Type {
	pair, ok := p.(pairVal)
	if !ok {
		return NewNil()
	}
	return pair.Cdr()
}

func Append(head Type, value Type) (out Type) {
	if head.Type() == types.Nil {
		return Cons(value, head)
	}
	switch head.Type() {
	case types.Pair:
		out = Cons(Car(head), Append(Cdr(head), value))
	default:
		out = Cons(head, value)
	}
	return out
}

func (pr pairVal) Car() Type {
	return pr.car
}

func (pr pairVal) Cdr() Type {
	return pr.cdr
}

func Reverse(input Type) (output Type) {
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

func (pr pairVal) Equal(p Type) bool {
	otherPair, ok := p.(Pair)
	if !ok {
		return false
	}

	a := Pair(pr)
	b := otherPair

	for {
		// compare current elements
		if !a.Car().Equal(b.Car()) {
			return false
		}

		aCdr := a.Cdr()
		bCdr := b.Cdr()

		aNext, aIsPair := aCdr.(Pair)
		bNext, bIsPair := bCdr.(Pair)

		// if neither cdr is a Pair (likely both Nil or immediate values), compare directly
		if !aIsPair && !bIsPair {
			return aCdr.Equal(bCdr)
		}

		// if one is a Pair and the other is not, lengths/structure differ
		if aIsPair != bIsPair {
			return false
		}

		// advance both lists
		a = aNext
		b = bNext
	}

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
