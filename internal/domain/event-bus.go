package domain

import "context"

type EventBus interface {
	Publish(context.Context, Event) error
}

type Event interface {
	Key() string
	Data() []byte
}
