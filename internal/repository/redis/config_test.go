package redisrepo

import (
	"testing"
	"time"
)

func TestRedisRepository_Enabled(t *testing.T) {
	tests := []struct {
		name     string
		useRedis bool
		want     bool
	}{
		{"enabled", true, true},
		{"disabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisRepository{UseRedis: tt.useRedis}
			if got := r.Enabled(); got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisRepository_GetTTL(t *testing.T) {
	tests := []struct {
		name string
		ttl  time.Duration
	}{
		{"1 minute", 1 * time.Minute},
		{"5 minutes", 5 * time.Minute},
		{"1 hour", 1 * time.Hour},
		{"zero", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisRepository{DefaultTTL: tt.ttl}
			if got := r.GetTTL(); got != tt.ttl {
				t.Errorf("GetTTL() = %v, want %v", got, tt.ttl)
			}
		})
	}
}
