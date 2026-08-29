package events

type Subscriber[T any] struct {
	ch chan T
}

func (s Subscriber[T]) Channel() chan T {
	return s.ch
}

func (s Subscriber[T]) Unsubscribe() {
	close(s.ch)
}

func NewSubscriber[T any]() Subscriber[T] {
	return Subscriber[T]{
		ch: make(chan T),
	}
}
