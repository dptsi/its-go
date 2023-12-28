package contracts

import (
	"context"
	"time"

	"bitbucket.org/dptsi/go-framework/web"
)

type SessionData interface {
	Id() string
	CSRFToken() string
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Clear()
	RegenerateId()
	RegenerateCSRFToken()
	Data() map[string]interface{}
	ExpiredAt() time.Time
}

type SessionStorage interface {
	Get(ctx context.Context, id string) (SessionData, error)
	Save(ctx context.Context, data SessionData) error
	Delete(ctx context.Context, id string) error
}

type SessionCookieWriter interface {
	Write(ctx *web.Context, sess SessionData)
}

type SessionService interface {
	Get(ctx *web.Context, key string) (interface{}, error)
	IsTokenMatch(ctx *web.Context, token string) (bool, error)
	Put(ctx *web.Context, key string, value interface{}) error
	Delete(ctx *web.Context, key string) error
	Clear(ctx *web.Context) error
	Regenerate(ctx *web.Context) error
	Invalidate(ctx *web.Context) error
	RegenerateToken(ctx *web.Context) error
	Start(ctx *web.Context) error
}
