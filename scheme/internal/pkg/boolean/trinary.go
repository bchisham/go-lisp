package boolean

func Trinary[T any](test bool, t T, f T) T {
	if test {
		return t
	}
	return f
}

func LazyTrinary[T any](test bool, t func() T, f func() T) T {
	if test {
		return t()
	}
	return f()
}
