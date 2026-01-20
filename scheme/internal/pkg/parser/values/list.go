package values

import "github.com/bchisham/go-lisp/scheme/internal/pkg/parser/types"

func NewList(v ...Interface) Interface {
	if len(v) == 0 {
		return Value{
			t:       types.List,
			ListVal: []Interface{},
		}
	}
	if len(v) == 1 && v[0].Type() == types.List {
		return v[0]
	}
	return Value{
		t:       types.List,
		ListVal: v,
	}
}
