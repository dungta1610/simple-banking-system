package ratelimit

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Limiter interface {
	IsAllowed(ctx context.Context, ip string, path string, limit int64, window time.Duration) bool
}

type RedisLimiter struct {
	rdb    *goredis.Client
	prefix string
}

func NewRedisLimiter(rdb *goredis.Client, prefix string) *RedisLimiter {
	if prefix == "" {
		prefix = "ratelimit"
	}

	return &RedisLimiter{rdb: rdb, prefix: prefix}
}

func (l *RedisLimiter) key(ip, path string) string {
	if ip == "" {
		ip = "unknown"
	}

	if path == "" {
		path = "global"
	}

	return fmt.Sprintf("%s:%s:%s", l.prefix, ip, path)
}

func (l *RedisLimiter) IsAllowed(ctx context.Context, ip string, path string, limit int64, window time.Duration) bool {
	if l == nil || l.rdb == nil {
		return true
	}

	if limit <= 0 {
		limit = 10
	}

	if window <= 0 {
		window = time.Minute
	}

	k := l.key(ip, path)
	count, err := l.rdb.Incr(ctx, k).Result()

	if err != nil {
		return true
	}

	if count == 1 {
		_ = l.rdb.Expire(ctx, k, window).Err()
	}

	return count <= limit
}
