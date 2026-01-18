package boolean

func NotFunc[T any](orig func(T) bool) func(T) bool {
	return func(e T) bool {
		return !orig(e)
	}
}
