package errors

func OrDie[T any](f func() (T, error)) T {
	v, err := f()
	if err != nil {
		panic(err)
	}
	return v
}

func Must[T any](f func() (T, error)) T {
	v, err := f()
	if err != nil {
		panic(err)
	}
	return v
}

func IgnoreError[T any](f func() (T, error)) T {
	v, _ := f()
	return v
}
