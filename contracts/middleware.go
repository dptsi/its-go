package contracts

import "bitbucket.org/dptsi/go-framework/web"

type Middleware interface {
	Handle(param interface{}) web.HandlerFunc
}

type MiddlewareService interface {
	Get(name string) (Middleware, error)
	Global() ([]Middleware, error)
}
