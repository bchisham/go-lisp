package values

import (
	"fmt"
	"math"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"
)

var Zero Numeric = Number{
	IntVal: 0,
	IsInt:  true,
}

var One Numeric = Number{
	IntVal: 1,
	IsInt:  true,
}

func NewInt(i int64) Interface {
	return Number{
		IsInt:  true,
		IntVal: i,
	}
}

func NewFloat(f float64) Interface {
	return Number{
		IsInt:    false,
		FloatVal: f,
	}
}

func IsNaN(n Numeric) bool {
	if n.IsFloat() {
		f, _ := n.AsFloat()
		return math.IsNaN(f)
	}
	return false
}

type Number struct {
	truthyValue
	FloatVal float64
	IntVal   int64
	IsInt    bool
}

func (n Number) Equal(p Interface) bool {
	if n.Type() != p.Type() {
		return false
	}
	o, ok := p.(Number)
	if !ok {
		return false
	}
	if n.IsInt != o.IsInt {
		return false
	}
	if n.IsInt {
		if n.IntVal != o.IntVal {
			return false
		}
	} else {
		if n.FloatVal != o.FloatVal {
			return false
		}
	}
	return true
}

func (n Number) Type() types.Type {
	if n.IsInt {
		return types.Int
	}
	return types.Float
}

func (n Number) IsTruthy() bool {
	return true
}

func (n Number) DisplayString() string {
	if n.IsInt {
		return fmt.Sprintf("%d", n.IntVal)
	}
	return fmt.Sprintf("%g", n.FloatVal)
}

func (n Number) WriteString() string {
	if n.IsInt {
		return fmt.Sprintf("%d", n.IntVal)
	}
	return fmt.Sprintf("%g", n.FloatVal)
}

func (n Number) IsInteger() bool {
	return n.IsInt
}

func (n Number) IsFloat() bool {
	return !n.IsInt
}

func (n Number) AsFloat() (float64, error) {
	if n.IsInt {
		return float64(n.IntVal), nil
	}
	return n.FloatVal, nil
}

func (n Number) AsInt() (int64, error) {
	if n.IsInt {
		return n.IntVal, nil
	}
	return 0, fmt.Errorf("cannot convert float to int")
}

func (n Number) Add(rhs Numeric) Numeric {
	if n.IsInt && rhs.IsInteger() {
		rhsInt, _ := rhs.AsInt()
		return Number{
			IntVal: n.IntVal + rhsInt,
			IsInt:  true,
		}
	}
	lhsFloat, _ := n.AsFloat()
	rhsFloat, _ := rhs.AsFloat()
	return Number{
		FloatVal: lhsFloat + rhsFloat,
		IsInt:    false,
	}
}

func (n Number) Sub(rhs Numeric) Numeric {
	if n.IsInt && rhs.IsInteger() {
		rhsInt, _ := rhs.AsInt()
		return Number{
			IntVal: n.IntVal - rhsInt,
			IsInt:  true,
		}
	}
	lhsFloat, _ := n.AsFloat()
	rhsFloat, _ := rhs.AsFloat()
	return Number{
		FloatVal: lhsFloat - rhsFloat,
		IsInt:    false,
	}
}

func (n Number) Mul(rhs Numeric) Numeric {
	if n.IsInt && rhs.IsInteger() {
		rhsInt, _ := rhs.AsInt()
		return Number{
			IntVal: n.IntVal * rhsInt,
			IsInt:  true,
		}
	}
	lhsFloat, _ := n.AsFloat()
	rhsFloat, _ := rhs.AsFloat()
	return Number{
		FloatVal: lhsFloat * rhsFloat,
		IsInt:    false,
	}
}

func (n Number) Div(rhs Numeric) (Numeric, error) {
	if n.IsInt && rhs.IsInteger() {
		rhsInt, _ := rhs.AsInt()
		if rhsInt == 0 {
			return nil, fmt.Errorf("divide by zero")
		}
		return Number{
			FloatVal: float64(n.IntVal) / float64(rhsInt),
			IsInt:    false,
		}, nil
	}
	lhsFloat, _ := n.AsFloat()
	rhsFloat, _ := rhs.AsFloat()
	if rhsFloat == 0 {
		return nil, fmt.Errorf("divide by zero")
	}
	return Number{
		FloatVal: lhsFloat / rhsFloat,
		IsInt:    false,
	}, nil
}

func (n Number) Mod(rhs Numeric) (Numeric, error) {
	if n.IsInt && rhs.IsInteger() {
		rhsInt, _ := rhs.AsInt()
		if rhsInt == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		return Number{
			IntVal: n.IntVal % rhsInt,
			IsInt:  true,
		}, nil
	}
	lhsFloat, _ := n.AsFloat()
	rhsFloat, _ := rhs.AsFloat()
	if rhsFloat == 0 {
		return nil, fmt.Errorf("modulo by zero")
	}
	return Number{
		FloatVal: float64(int(lhsFloat) % int(rhsFloat)),
		IsInt:    false,
	}, nil
}

func (n Number) Neg() Numeric {
	if n.IsInt {
		return Number{
			IntVal: -n.IntVal,
			IsInt:  true,
		}
	}
	return Number{
		FloatVal: -n.FloatVal,
		IsInt:    false,
	}
}

func (n Number) Abs() Numeric {
	if n.IsInt {
		if n.IntVal < 0 {
			return Number{
				IntVal: -n.IntVal,
				IsInt:  true,
			}
		}
		return n
	}
	if n.FloatVal < 0 {
		return Number{
			FloatVal: -n.FloatVal,
			IsInt:    false,
		}
	}
	return n
}

func (n Number) GreaterThan(rhs Numeric) bool {
	if n.IsInt && rhs.IsInteger() {
		rhsInt, _ := rhs.AsInt()
		return n.IntVal > rhsInt
	}
	rhsFloat, err := n.AsFloat()
	if err != nil {
		return false
	}
	return n.FloatVal > rhsFloat
}
func (n Number) LessThan(rhs Numeric) bool {
	if n.IsInt && rhs.IsInteger() {
		rhsInt, _ := rhs.AsInt()
		return n.IntVal < rhsInt
	}
	rhsFloat, err := rhs.AsFloat()
	if err != nil {
		return false
	}
	return n.FloatVal < rhsFloat
}

func (n Number) GreaterThanOrEqual(rhs Numeric) bool {
	if n.IsInt && rhs.IsInteger() {
		rhsInt, _ := rhs.AsInt()
		return n.IntVal >= rhsInt
	}
	rhsFloat, err := rhs.AsFloat()
	if err != nil {
		return false
	}
	return n.FloatVal >= rhsFloat
}
func (n Number) LessThanOrEqual(rhs Numeric) bool {
	if n.IsInt && rhs.IsInteger() {
		rhsInt, _ := rhs.AsInt()
		return n.IntVal <= rhsInt
	}
	rhsFloat, err := rhs.AsFloat()
	if err != nil {
		return false
	}
	return n.FloatVal <= rhsFloat
}
