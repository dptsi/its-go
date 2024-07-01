package event

import (
	"context"
	"fmt"
	"log"

	"github.com/dptsi/its-go/app"
	"github.com/dptsi/its-go/contracts"
)

type Service struct {
	application    contracts.Application
	listenersCount map[string]int
}

func NewService(application contracts.Application) *Service {
	return &Service{
		application:    application,
		listenersCount: make(map[string]int),
	}
}

func (s *Service) Dispatch(ctx context.Context, payload contracts.Event) {
	name := payload.Name()
	listenerCount, exists := s.listenersCount[name]
	if !exists {
		// log.Printf("event service: dispatch: event name %s is not found", name)
		return
	}

	for i := 0; i < listenerCount; i += 1 {
		listener, err := app.Make[contracts.EventListener](s.application, s.getListenerServiceKey(name, i))
		if err != nil {
			log.Printf("event service: dispatch: error when constructing listener for event %s with index %d: %s", name, i, err.Error())
			log.Printf("event service: dispatch: remaining listener for event %s will not be executed", name)
			return
		}

		if err := listener.Handle(ctx, payload); err != nil {
			log.Printf("event service: dispatch: listener %s return error: %s", listener.Name(), err.Error())
			log.Printf("event service: dispatch: remaining listener for event %s will not be executed", name)
			return
		}
	}
}

func (s *Service) Register(name string, listenersConstructor []contracts.EventListenerConstructor) {
	if _, exists := s.listenersCount[name]; exists {
		log.Fatalf("event service: dispatch: event with name %s already exist", name)
	}
	s.listenersCount[name] = len(listenersConstructor)

	for i, constructor := range listenersConstructor {
		app.Bind(
			s.application,
			s.getListenerServiceKey(name, i),
			constructor,
		)
	}
}

func (s *Service) getListenerServiceKey(name string, index int) string {
	return fmt.Sprintf("event.%s.listeners.%d", name, index)
}
