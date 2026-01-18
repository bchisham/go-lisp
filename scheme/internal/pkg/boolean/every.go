package boolean

func EveryFuc[T any](f ...func(T) bool) func(t T) bool {
	return func(e T) bool {
		for _, fn := range f {
			if !fn(e) {
				return false
			}
		}
		return true
	}
}
