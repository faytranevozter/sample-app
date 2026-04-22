package connection

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedis connects to Redis database using go-redis and returns the client.
// It takes a timeout duration and a database URL as parameters.
// The database URL should be in the format: redis://username:password@host:port/db?query_params
//
// Example: redis://default:password@localhost:6379/0?ssl=false
func NewRedis(timeout time.Duration, dbURL string) *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	opts, err := redis.ParseURL(dbURL)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opts)

	if ping := client.Ping(ctx); ping.Err() != nil {
		panic(ping.Err())
	}

	return client
}
