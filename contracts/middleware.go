package contracts

import "bitbucket.org/dptsi/its-go/web"

type Middleware interface {
	Handle(param interface{}) web.HandlerFunc
}

type MiddlewareConstructor = func(application Application) (Middleware, error)
type MiddlewareService interface {
	Use(name string, params interface{}) web.HandlerFunc
	Global() []web.HandlerFunc
	Register(name string, constructor MiddlewareConstructor) error
}
