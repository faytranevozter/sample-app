package redisrepo

import (
	"context"

	"github.com/sirupsen/logrus"
)

func (r *RedisRepository) Get(ctx context.Context, key string) (value []byte, err error) {
	if value, err = r.Conn.Get(ctx, r.Prefix+key).Bytes(); err != nil {
		logrus.Error("Redis Get:", err)
		return
	}

	return
}
