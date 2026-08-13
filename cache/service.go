package cache

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/dptsi/its-go/contracts"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultKeyPrefix = "cache:"
)

var (
	ErrCacheMiss = errors.New("cache: key not found")
)

func formatKey(key string) string {
	if strings.HasPrefix(key, DefaultKeyPrefix) {
		return key
	}
	return DefaultKeyPrefix + key
}

func formatKeys(keys ...string) []string {
	formatted := make([]string, len(keys))
	for i, k := range keys {
		formatted[i] = formatKey(k)
	}
	return formatted
}

type Service struct {
	client *redis.Client
}

func NewService(client *redis.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Client() *redis.Client {
	return s.client
}

func (s *Service) IsAvailable() bool {
	return s.client != nil
}

func (s *Service) Get(ctx context.Context, key string, dest interface{}) error {
	if s.client == nil {
		return ErrCacheMiss
	}

	val, err := s.client.Get(ctx, formatKey(key)).Bytes()
	if errors.Is(err, redis.Nil) {
		return ErrCacheMiss
	}
	if err != nil {
		return err
	}

	switch target := dest.(type) {
	case *string:
		*target = string(val)
		return nil
	case *[]byte:
		*target = val
		return nil
	default:
		return json.Unmarshal(val, dest)
	}
}

func (s *Service) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if s.client == nil {
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		var err error
		data, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}

	return s.client.Set(ctx, formatKey(key), data, ttl).Err()
}

func (s *Service) Delete(ctx context.Context, keys ...string) error {
	if s.client == nil || len(keys) == 0 {
		return nil
	}
	return s.client.Del(ctx, formatKeys(keys...)...).Err()
}

func (s *Service) DeletePrefix(ctx context.Context, prefix string) error {
	if s.client == nil || prefix == "" {
		return nil
	}

	pattern := formatKey(prefix)
	if !strings.HasSuffix(pattern, "*") {
		pattern = pattern + "*"
	}

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = s.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := s.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (s *Service) Has(ctx context.Context, key string) bool {
	if s.client == nil {
		return false
	}
	count, err := s.client.Exists(ctx, formatKey(key)).Result()
	return err == nil && count > 0
}

func (s *Service) Remember(ctx context.Context, key string, ttl time.Duration, dest interface{}, fn func() (interface{}, error)) error {
	if s.client != nil && dest != nil {
		err := s.Get(ctx, key, dest)
		if err == nil {
			return nil
		}
	}

	val, err := fn()
	if err != nil {
		return err
	}

	if dest != nil {
		bytes, err := json.Marshal(val)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(bytes, dest); err != nil {
			return err
		}
	}

	if s.client != nil {
		_ = s.Set(ctx, key, val, ttl)
	}

	return nil
}

func (s *Service) Flush(ctx context.Context) error {
	if s.client == nil {
		return nil
	}
	return s.client.FlushDB(ctx).Err()
}

// Atomic Counters & Key Expiration

func (s *Service) Increment(ctx context.Context, key string) (int64, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.Incr(ctx, formatKey(key)).Result()
}

func (s *Service) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.IncrBy(ctx, formatKey(key), value).Result()
}

func (s *Service) Decrement(ctx context.Context, key string) (int64, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.Decr(ctx, formatKey(key)).Result()
}

func (s *Service) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.DecrBy(ctx, formatKey(key), value).Result()
}

func (s *Service) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if s.client == nil {
		return nil
	}
	return s.client.Expire(ctx, formatKey(key), ttl).Err()
}

func (s *Service) TTL(ctx context.Context, key string) (time.Duration, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.TTL(ctx, formatKey(key)).Result()
}

// Geospatial Operations

func (s *Service) GeoAdd(ctx context.Context, key string, locations ...*redis.GeoLocation) error {
	if s.client == nil || len(locations) == 0 {
		return nil
	}
	return s.client.GeoAdd(ctx, formatKey(key), locations...).Err()
}

func (s *Service) GeoDist(ctx context.Context, key string, member1, member2, unit string) (float64, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.GeoDist(ctx, formatKey(key), member1, member2, unit).Result()
}

func (s *Service) GeoPos(ctx context.Context, key string, members ...string) ([]*redis.GeoPos, error) {
	if s.client == nil || len(members) == 0 {
		return nil, nil
	}
	return s.client.GeoPos(ctx, formatKey(key), members...).Result()
}

func (s *Service) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	if s.client == nil {
		return nil, nil
	}
	return s.client.GeoRadius(ctx, formatKey(key), longitude, latitude, query).Result()
}

func (s *Service) GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) ([]string, error) {
	if s.client == nil {
		return nil, nil
	}
	return s.client.GeoSearch(ctx, formatKey(key), q).Result()
}

func (s *Service) GeoSearchLocation(ctx context.Context, key string, q *redis.GeoSearchLocationQuery) ([]redis.GeoLocation, error) {
	if s.client == nil {
		return nil, nil
	}
	return s.client.GeoSearchLocation(ctx, formatKey(key), q).Result()
}

var defaultService contracts.CacheService

// SetDefault sets the default CacheService used by package-level helpers like cache.Remember.
func SetDefault(s contracts.CacheService) {
	defaultService = s
}

// Default returns the default CacheService.
func Default() contracts.CacheService {
	return defaultService
}

// Remember is a type-safe generic helper that uses the default CacheService to cache the result of a function.
func Remember[T any](ctx context.Context, key string, ttl time.Duration, fn func() (T, error)) (T, error) {
	return RememberWith[T](ctx, defaultService, key, ttl, fn)
}

// RememberWith allows explicitly specifying a custom CacheService (e.g. for secondary Redis connections).
func RememberWith[T any](ctx context.Context, s contracts.CacheService, key string, ttl time.Duration, fn func() (T, error)) (T, error) {
	if s != nil {
		var cached T
		err := s.Get(ctx, key, &cached)
		if err == nil {
			return cached, nil
		}
	}

	val, err := fn()
	if err != nil {
		var zero T
		return zero, err
	}

	if s != nil {
		_ = s.Set(ctx, key, val, ttl)
	}

	return val, nil
}
