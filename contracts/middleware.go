package contracts

import "bitbucket.org/dptsi/go-framework/web"

type Middleware interface {
	Handle(param interface{}) web.HandlerFunc
}

type MiddlewareService interface {
	Use(name string, params interface{}) web.HandlerFunc
	Global() web.HandlerFunc
}
