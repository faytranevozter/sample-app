package middleware

import (
	jwt_helper "app/helpers/jsonwebtoken"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Middleware interface {
	Auth() gin.HandlerFunc
	Cors() gin.HandlerFunc
	Logger(writer io.Writer) gin.HandlerFunc
	Recovery() gin.HandlerFunc
	Cache(expiry ...time.Duration) gin.HandlerFunc
}

type appMiddleware struct {
	secret string
	cache  CacheConfig
}

type CacheConfig struct {
	enabled     bool
	headerKeys  []string
	store       *redis.Client
	storeTTL    time.Duration
	cachePrefix string
}

func NewMiddleware(redis *redis.Client) Middleware {
	ttl, _ := time.ParseDuration(os.Getenv("REDIS_TTL"))
	// default ttl redis
	if ttl == 0 {
		ttl = 1 * time.Minute
	}

	useRedis, _ := strconv.ParseBool(os.Getenv("USE_REDIS"))
	redisKeyPrefix := os.Getenv("REDIS_KEY_PREFIX")

	return &appMiddleware{
		secret: jwt_helper.GetJwtCredential().Member.Secret,
		cache: CacheConfig{
			enabled:     useRedis,
			store:       redis,
			cachePrefix: redisKeyPrefix + "gin:",
			storeTTL:    ttl,
			headerKeys: []string{
				"User-Agent",
				"Accept",
				"Accept-Encoding",
				"Accept-Language",
				"Cookie",
			},
		},
	}
}
