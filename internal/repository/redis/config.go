package redisrepo

import (
	"time"
)

func (r *RedisRepository) Enabled() bool {
	return r.UseRedis
}

func (r *RedisRepository) GetTTL() time.Duration {
	return r.DefaultTTL
}
