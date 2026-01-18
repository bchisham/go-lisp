package list

func Apply[List ~[]T, T any, U any](l List, transform func(T) U) []U {
	result := make([]U, 0, len(l))
	for _, item := range l {
		result = append(result, transform(item))
	}
	return result
}
