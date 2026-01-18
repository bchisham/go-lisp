package list

func New[T any](t ...T) []T {
	return t
}

func Cons[T any](head T, tail []T) []T {
	return append(New(head), tail...)
}

func Car[T any](list []T) T {
	var zero T
	if len(list) == 0 {
		return zero
	}
	return list[0]
}

func Cdr[T any](list []T) []T {
	if len(list) == 0 {
		return nil
	}
	return list[1:]
}
