package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/dptsi/its-go/contracts"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

const (
	DefaultKeyPrefix = "cache:"
	tagPrefix        = "cache:tag:"
)

var (
	ErrCacheMiss = errors.New("cache: key not found")
	sfGroup      singleflight.Group
)

// Config holds cache configuration options.
type Config struct {
	// Prefix is the key prefix prepended to all cache keys. Defaults to "cache:".
	Prefix string
}

func formatKeyWithPrefix(prefix, key string) string {
	prefix = normalizePrefix(prefix)
	if strings.HasPrefix(key, prefix) {
		return key
	}
	return prefix + key
}

// normalizePrefix returns the default prefix when empty and guarantees the
// prefix ends with ":" so keys are always separated consistently.
func normalizePrefix(prefix string) string {
	if prefix == "" {
		return DefaultKeyPrefix
	}
	if !strings.HasSuffix(prefix, ":") {
		prefix += ":"
	}
	return prefix
}

type Service struct {
	client *redis.Client
	prefix string
	tags   map[string][]string // key -> tags (in-memory tracking for tagged keys)
	mu     sync.RWMutex
}

// NewService creates a new cache service with the default "cache:" prefix.
func NewService(client *redis.Client) *Service {
	return &Service{
		client: client,
		prefix: DefaultKeyPrefix,
		tags:   make(map[string][]string),
	}
}

// NewServiceWithConfig creates a cache service with custom configuration.
// If cfg.Prefix does not end with ":", it is appended automatically.
func NewServiceWithConfig(client *redis.Client, cfg Config) *Service {
	return &Service{
		client: client,
		prefix: normalizePrefix(cfg.Prefix),
		tags:   make(map[string][]string),
	}
}

func (s *Service) Client() *redis.Client {
	return s.client
}

func (s *Service) IsAvailable() bool {
	return s.client != nil
}

func (s *Service) formatKey(key string) string {
	return formatKeyWithPrefix(s.prefix, key)
}

func (s *Service) Get(ctx context.Context, key string, dest interface{}) error {
	if s.client == nil {
		return ErrCacheMiss
	}

	val, err := s.client.Get(ctx, s.formatKey(key)).Bytes()
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

	data, err := s.marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, s.formatKey(key), data, ttl).Err()
}

// Add stores the value only if the key does not already exist (Laravel Cache::add).
// Returns true if the value was stored, false if the key already existed.
func (s *Service) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	if s.client == nil {
		return false, nil
	}

	data, err := s.marshal(value)
	if err != nil {
		return false, err
	}

	ok, err := s.client.SetNX(ctx, s.formatKey(key), data, ttl).Result()
	return ok, err
}

// Pull retrieves a value and then deletes it (Laravel Cache::pull).
// Returns ErrCacheMiss if the key does not exist.
func (s *Service) Pull(ctx context.Context, key string, dest interface{}) error {
	if err := s.Get(ctx, key, dest); err != nil {
		return err
	}
	return s.Forget(ctx, key)
}

// Forever stores a value without expiration (Laravel Cache::forever).
func (s *Service) Forever(ctx context.Context, key string, value interface{}) error {
	return s.Set(ctx, key, value, 0)
}

// RememberForever retrieves a value or stores it forever using the callback (Laravel Cache::rememberForever).
func (s *Service) RememberForever(ctx context.Context, key string, dest interface{}, fn func() (interface{}, error)) error {
	if s.client != nil && dest != nil {
		if err := s.Get(ctx, key, dest); err == nil {
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
		_ = s.Forever(ctx, key, val)
	}

	return nil
}

func (s *Service) marshal(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return json.Marshal(v)
	}
}

func (s *Service) Forget(ctx context.Context, pattern string) error {
	if s.client == nil || pattern == "" {
		return nil
	}

	if strings.Contains(pattern, "*") {
		formattedPattern := s.formatKey(pattern)
		if !strings.HasSuffix(formattedPattern, "*") {
			formattedPattern += "*"
		}

		var cursor uint64
		for {
			var keys []string
			var err error
			keys, cursor, err = s.client.Scan(ctx, cursor, formattedPattern, 100).Result()
			if err != nil {
				log.Printf("[cache] Forget SCAN error for pattern %q: %v", formattedPattern, err)
				return nil
			}

			if len(keys) > 0 {
				if err := s.client.Del(ctx, keys...).Err(); err != nil {
					log.Printf("[cache] Forget DEL error for pattern %q: %v", formattedPattern, err)
					return nil
				}
			}

			if cursor == 0 {
				break
			}
		}
		return nil
	}

	return s.client.Del(ctx, s.formatKey(pattern)).Err()
}

func (s *Service) Has(ctx context.Context, key string) bool {
	if s.client == nil {
		return false
	}
	count, err := s.client.Exists(ctx, s.formatKey(key)).Result()
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

// Flush removes all keys managed by this cache service (scoped by prefix).
// Unlike a raw FLUSHDB, this only deletes keys matching the service prefix,
// making it safe on shared Redis instances.
func (s *Service) Flush(ctx context.Context) error {
	if s.client == nil {
		return nil
	}
	return s.Forget(ctx, "*")
}

// Atomic Counters & Key Expiration

func (s *Service) Increment(ctx context.Context, key string) (int64, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.Incr(ctx, s.formatKey(key)).Result()
}

func (s *Service) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.IncrBy(ctx, s.formatKey(key), value).Result()
}

func (s *Service) Decrement(ctx context.Context, key string) (int64, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.Decr(ctx, s.formatKey(key)).Result()
}

func (s *Service) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.DecrBy(ctx, s.formatKey(key), value).Result()
}

func (s *Service) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if s.client == nil {
		return nil
	}
	return s.client.Expire(ctx, s.formatKey(key), ttl).Err()
}

func (s *Service) TTL(ctx context.Context, key string) (time.Duration, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.TTL(ctx, s.formatKey(key)).Result()
}

// Geospatial Operations

func (s *Service) GeoAdd(ctx context.Context, key string, locations ...*redis.GeoLocation) error {
	if s.client == nil || len(locations) == 0 {
		return nil
	}
	return s.client.GeoAdd(ctx, s.formatKey(key), locations...).Err()
}

func (s *Service) GeoDist(ctx context.Context, key string, member1, member2, unit string) (float64, error) {
	if s.client == nil {
		return 0, nil
	}
	return s.client.GeoDist(ctx, s.formatKey(key), member1, member2, unit).Result()
}

func (s *Service) GeoPos(ctx context.Context, key string, members ...string) ([]*redis.GeoPos, error) {
	if s.client == nil || len(members) == 0 {
		return nil, nil
	}
	return s.client.GeoPos(ctx, s.formatKey(key), members...).Result()
}

func (s *Service) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	if s.client == nil {
		return nil, nil
	}
	return s.client.GeoRadius(ctx, s.formatKey(key), longitude, latitude, query).Result()
}

func (s *Service) GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) ([]string, error) {
	if s.client == nil {
		return nil, nil
	}
	return s.client.GeoSearch(ctx, s.formatKey(key), q).Result()
}

func (s *Service) GeoSearchLocation(ctx context.Context, key string, q *redis.GeoSearchLocationQuery) ([]redis.GeoLocation, error) {
	if s.client == nil {
		return nil, nil
	}
	return s.client.GeoSearchLocation(ctx, s.formatKey(key), q).Result()
}

// --- Tags (Laravel Cache::tags) ---

// TaggedSet stores a value associated with tags. Tagged keys can later be
// invalidated together via ForgetTags without scanning the entire keyspace.
func (s *Service) TaggedSet(ctx context.Context, tags []string, key string, value interface{}, ttl time.Duration) error {
	if err := s.Set(ctx, key, value, ttl); err != nil {
		return err
	}
	if s.client == nil || len(tags) == 0 {
		return nil
	}

	formattedKey := s.formatKey(key)
	for _, tag := range tags {
		tagKey := tagPrefix + tag
		if err := s.client.SAdd(ctx, tagKey, formattedKey).Err(); err != nil {
			log.Printf("[cache] TaggedSet SADD error for tag %q: %v", tag, err)
		}
	}

	s.mu.Lock()
	s.tags[formattedKey] = tags
	s.mu.Unlock()

	return nil
}

// ForgetTags removes all keys associated with the given tags (Laravel Cache::tags(...)->flush()).
func (s *Service) ForgetTags(ctx context.Context, tags ...string) error {
	if s.client == nil || len(tags) == 0 {
		return nil
	}

	for _, tag := range tags {
		tagKey := tagPrefix + tag
		keys, err := s.client.SMembers(ctx, tagKey).Result()
		if err != nil {
			log.Printf("[cache] ForgetTags SMEMBERS error for tag %q: %v", tag, err)
			continue
		}

		if len(keys) > 0 {
			if err := s.client.Del(ctx, keys...).Err(); err != nil {
				log.Printf("[cache] ForgetTags DEL error for tag %q: %v", tag, err)
			}
		}

		_ = s.client.Del(ctx, tagKey).Err()
	}

	return nil
}

// --- Flexible (stale-while-revalidate, Laravel 11+ style) ---

// Flexible returns the cached value immediately, even if stale. If the value
// is older than `ttl` but within `grace`, it triggers a background refresh.
// If the value is beyond grace or missing, fn is called synchronously.
func (s *Service) Flexible(ctx context.Context, key string, ttl, grace time.Duration, dest interface{}, fn func() (interface{}, error)) error {
	if s.client != nil && dest != nil {
		if err := s.Get(ctx, key, dest); err == nil {
			// Check remaining TTL to decide if background refresh is needed
			remaining, _ := s.TTL(ctx, key)
			if remaining > 0 && remaining < grace {
				// Stale but within grace: refresh in background
				go func() {
					val, err := fn()
					if err != nil {
						log.Printf("[cache] Flexible background refresh error for key %q: %v", key, err)
						return
					}
					_ = s.Set(context.Background(), key, val, ttl)
				}()
			}
			return nil
		}
	}

	// Cache miss or beyond grace: synchronous fetch
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

var defaultService contracts.CacheService
var namedServices = map[string]contracts.CacheService{}

func Register(name string, s contracts.CacheService) {
	if name == "" {
		return
	}
	namedServices[name] = s
}

func Get(name string) contracts.CacheService {
	if name == "" {
		return defaultService
	}
	if svc, ok := namedServices[name]; ok {
		return svc
	}
	return defaultService
}

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
	return RememberOn[T](ctx, "", key, ttl, fn)
}

// RememberOn uses a named cache service when provided, otherwise falls back to the default service.
// It includes singleflight stampede protection: concurrent callers for the same
// key will share a single underlying fetch instead of all hitting the source.
func RememberOn[T any](ctx context.Context, redisName string, key string, ttl time.Duration, fn func() (T, error)) (T, error) {
	svc := Get(redisName)
	if svc == nil {
		return fn()
	}

	var cached T
	if err := svc.Get(ctx, key, &cached); err == nil {
		return cached, nil
	}

	// Singleflight: dedupe concurrent misses for the same key.
	flightKey := redisName + "\x00" + key
	v, err, _ := sfGroup.Do(flightKey, func() (interface{}, error) {
		// Double-check cache in case another goroutine populated it
		// between our first Get and acquiring the flight.
		var fresh T
		if err := svc.Get(ctx, key, &fresh); err == nil {
			return fresh, nil
		}

		val, err := fn()
		if err != nil {
			var zero T
			return zero, err
		}

		_ = svc.Set(ctx, key, val, ttl)
		return val, nil
	})
	if err != nil {
		var zero T
		return zero, err
	}
	return v.(T), nil
}

// FlexibleOn is the generic, named-service variant of Service.Flexible.
// It serves the cached value immediately, refreshing in the background when
// stale but within grace, and fetching synchronously on miss or beyond grace.
func FlexibleOn[T any](ctx context.Context, redisName string, key string, ttl, grace time.Duration, fn func() (T, error)) (T, error) {
	svc := Get(redisName)
	if svc == nil {
		return fn()
	}

	var cached T
	if err := svc.Get(ctx, key, &cached); err == nil {
		remaining, _ := svc.TTL(ctx, key)
		if remaining > 0 && remaining < grace {
			// Stale but within grace: refresh in background.
			go func() {
				val, err := fn()
				if err != nil {
					log.Printf("[cache] FlexibleOn background refresh error for key %q: %v", key, err)
					return
				}
				_ = svc.Set(context.Background(), key, val, ttl)
			}()
		}
		return cached, nil
	}

	// Cache miss or beyond grace: synchronous fetch.
	val, err := fn()
	if err != nil {
		var zero T
		return zero, err
	}

	_ = svc.Set(ctx, key, val, ttl)
	return val, nil
}

// Flexible is the generic variant using the default cache service.
func Flexible[T any](ctx context.Context, key string, ttl, grace time.Duration, fn func() (T, error)) (T, error) {
	return FlexibleOn[T](ctx, "", key, ttl, grace, fn)
}

// RememberWith keeps backward compatibility with older explicit service-based code.
func RememberWith[T any](ctx context.Context, s contracts.CacheService, key string, ttl time.Duration, fn func() (T, error)) (T, error) {
	if s == nil {
		return fn()
	}

	var cached T
	if err := s.Get(ctx, key, &cached); err == nil {
		return cached, nil
	}

	val, err := fn()
	if err != nil {
		var zero T
		return zero, err
	}

	_ = s.Set(ctx, key, val, ttl)
	return val, nil
}

func Forget(ctx context.Context, pattern string) error {
	return ForgetOn(ctx, "", pattern)
}

func ForgetOn(ctx context.Context, redisName string, pattern string) error {
	svc := Get(redisName)
	if svc == nil {
		return nil
	}
	return svc.Forget(ctx, pattern)
}
