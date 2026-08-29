package events

import "context"

type Publisher[T any] struct {
	ctx         context.Context
	subscribers []chan T
}
