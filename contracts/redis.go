package contracts

import "github.com/dptsi/its-go/redis"

type RedisService interface {
	GetClient(name string) *redis.Client
	GetDefault() *redis.Client
}
