package redisrepo

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func (r *RedisRepository) Set(ctx context.Context, key string, value []byte, expiration *time.Duration) (err error) {
	if expiration == nil {
		defaultTTL := r.GetTTL()
		expiration = &defaultTTL
	}
	if res := r.Conn.Set(ctx, r.Prefix+key, value, *expiration); res.Err() != nil {
		logrus.Error("Redis Set:", err)
		return
	}

	return
}
