package redis

import "github.com/redis/go-redis/v9"

type (
	Client                 = redis.Client
	GeoLocation            = redis.GeoLocation
	GeoPos                 = redis.GeoPos
	GeoRadiusQuery         = redis.GeoRadiusQuery
	GeoSearchQuery         = redis.GeoSearchQuery
	GeoSearchLocationQuery = redis.GeoSearchLocationQuery
	Options                = redis.Options
)

var (
	Nil = redis.Nil
)
