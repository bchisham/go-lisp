package values

func NewList(v ...Interface) Interface {
	if len(v) == 0 {
		return NewNil()
	}

	var list = NewNil()

	for i := len(v) - 1; i >= 0; i-- {
		list = Cons(v[i], list)
	}
	return list
}
