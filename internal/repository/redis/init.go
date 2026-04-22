package redisrepo

import (
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	Conn       *redis.Client
	Prefix     string
	UseRedis   bool
	DefaultTTL time.Duration
}

func NewRedisRepo(Conn *redis.Client) *RedisRepository {
	ttl, _ := time.ParseDuration(os.Getenv("REDIS_TTL"))
	// default ttl redis
	if ttl == 0 {
		ttl = 1 * time.Minute
	}

	useRedis, _ := strconv.ParseBool(os.Getenv("USE_REDIS"))
	redisKeyPrefix := os.Getenv("REDIS_KEY_PREFIX")

	return &RedisRepository{
		Conn:       Conn,
		Prefix:     redisKeyPrefix + "appc:",
		UseRedis:   useRedis,
		DefaultTTL: ttl,
	}
}
