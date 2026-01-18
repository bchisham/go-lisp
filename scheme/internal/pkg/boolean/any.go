package boolean

func AnyFunc[T any](f ...func(T) bool) func(T) bool {
	return func(t T) bool {
		var result bool
		for _, x := range f {
			if x(t) {
				result = true
				goto done
			}
		}
	done:
		return result
	}
}
