package types

import "errors"

type NativeFunc string

var ErrFuncNotDefined = errors.New("native function not defined")

const (
	Format  NativeFunc = "format"
	Write   NativeFunc = "write"
	Display NativeFunc = "display"
	//Quot    NativeFunc = "quot"
)

func FromString(s string) (NativeFunc, error) {
	n := NativeFunc(s)
	if n != "" {
		return n, nil
	}
	return "", ErrFuncNotDefined
}

func (n NativeFunc) String() string { return string(n) }
