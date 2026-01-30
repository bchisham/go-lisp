package values

type truthyValue struct{}

func (t truthyValue) IsTruthy() bool {
	return true
}
