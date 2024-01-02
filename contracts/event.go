package contracts

import (
	"context"
	"time"

	"bitbucket.org/dptsi/go-framework/app"
)

// The minimal interface DomainEvent, implemented by all events, ensures support of an occurredOn() accessor.
// It enforces a basic contract for all events.
//
// References:
//  1. Implementing Domain-Driven Design, Vaughn Vernon
type Event interface {
	OccuredOn() time.Time
	JSON() ([]byte, error)
}

type EventListener interface {
	Name() string
	Handle(ctx context.Context, event Event) error
}

type EventListenerConstructor = func(application *app.Application) (EventListener, error)

type EventService interface {
	Dispatch(ctx context.Context, name string, payload Event)
	Register(name string, listenersConstructor []EventListenerConstructor)
}
