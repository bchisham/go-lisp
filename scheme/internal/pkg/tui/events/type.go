package events

import "time"

type EventName string

type Event interface {
	GetName() EventName
	HappenedAt() time.Time
}

type Type struct {
	At   time.Time
	Name EventName
}

func (t Type) GetName() EventName {
	return t.Name
}

func (t Type) HappenedAt() time.Time {
	return t.At
}
