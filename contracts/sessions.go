package contracts

import (
	"context"
	"time"
)

type SessionData interface {
	Id() string
	CSRFToken() string
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Clear()
	RegenerateId()
	Invalidate()
	RegenerateCSRFToken()
	Data() map[string]interface{}
	ExpiredAt() time.Time
}

type SessionStorage interface {
	Get(ctx context.Context, id string) (SessionData, error)
	Save(ctx context.Context, data SessionData) error
	Delete(ctx context.Context, id string) error
}
