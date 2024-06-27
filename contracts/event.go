package contracts

import (
	"context"
	"time"
)

// The minimal interface DomainEvent, implemented by all events, ensures support of an occurredOn() accessor.
// It enforces a basic contract for all events.
//
// References:
//  1. Implementing Domain-Driven Design, Vaughn Vernon
type Event interface {
	Name() string
	OccuredOn() time.Time
	JSON() ([]byte, error)
}

type EventListener interface {
	Name() string
	Handle(ctx context.Context, event Event) error
}

type EventListenerConstructor = func(application Application) (EventListener, error)

type EventService interface {
	Dispatch(ctx context.Context, payload Event)
	Register(name string, listenersConstructor []EventListenerConstructor)
}
