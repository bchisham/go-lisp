package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

func NewList(v ...Value) Value {
	if len(v) == 0 {
		return Value{
			Type:    types.List,
			ListVal: []Value{},
		}
	}
	if len(v) == 1 && v[0].Type == types.List {
		return v[0]
	}
	return Value{
		Type:    types.List,
		ListVal: v,
	}
}
