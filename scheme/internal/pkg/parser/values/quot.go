package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

type Quot interface {
	Interface
	GetValue() Interface
}

type quotValue struct {
	value Interface
}

func NewQuot(value Interface) Quot {
	return quotValue{value: value}
}

func (q quotValue) IsTruthy() bool {
	return true
}

func (q quotValue) DisplayString() string {
	return q.value.DisplayString()
}

func (q quotValue) WriteString() string {
	return q.value.WriteString()
}

func NewQuotValue(value Interface) Interface {
	return quotValue{value: value}
}

func (q quotValue) GetValue() Interface {
	return q.value
}
func (q quotValue) Equal(p Interface) bool {
	other, ok := p.(quotValue)
	if !ok {
		return false
	}
	return q.value.Equal(other.value)
}

func (q quotValue) Type() types.Type {
	return types.Quot
}
