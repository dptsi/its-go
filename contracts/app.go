package contracts

import "time"

// The minimal interface DomainEvent, implemented by all events, ensures support of an occurredOn() accessor.
// It enforces a basic contract for all events.
//
// References:
//  1. Implementing Domain-Driven Design, Vaughn Vernon
type Event interface {
	OccuredOn() time.Time
	JSON() ([]byte, error)
}
