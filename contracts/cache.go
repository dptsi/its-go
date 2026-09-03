package contracts

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	// Core Cache Operations
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Forget(ctx context.Context, pattern string) error
	Has(ctx context.Context, key string) bool
	Remember(ctx context.Context, key string, ttl time.Duration, dest interface{}, fn func() (interface{}, error)) error
	Flush(ctx context.Context) error
	Client() *redis.Client

	// Laravel-inspired Operations
	Add(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)
	Pull(ctx context.Context, key string, dest interface{}) error
	Forever(ctx context.Context, key string, value interface{}) error
	RememberForever(ctx context.Context, key string, dest interface{}, fn func() (interface{}, error)) error

	// Tags (Laravel Cache::tags)
	TaggedSet(ctx context.Context, tags []string, key string, value interface{}, ttl time.Duration) error
	ForgetTags(ctx context.Context, tags ...string) error

	// Flexible (stale-while-revalidate)
	Flexible(ctx context.Context, key string, ttl, grace time.Duration, dest interface{}, fn func() (interface{}, error)) error

	// Atomic Counters & Key Expiration
	Increment(ctx context.Context, key string) (int64, error)
	IncrementBy(ctx context.Context, key string, value int64) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)
	DecrementBy(ctx context.Context, key string, value int64) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Geospatial Operations
	GeoAdd(ctx context.Context, key string, locations ...*redis.GeoLocation) error
	GeoDist(ctx context.Context, key string, member1, member2, unit string) (float64, error)
	GeoPos(ctx context.Context, key string, members ...string) ([]*redis.GeoPos, error)
	GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error)
	GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) ([]string, error)
	GeoSearchLocation(ctx context.Context, key string, q *redis.GeoSearchLocationQuery) ([]redis.GeoLocation, error)
}
