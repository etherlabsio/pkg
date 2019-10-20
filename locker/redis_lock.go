package locker

import (
	"context"

	"github.com/bsm/redislock"
)

// NewRedisLocker creates a Locker implementation for Redis
func NewRedisLocker(c redislock.RedisClient) *RedisLocker {
	return &RedisLocker{
		locker: redislock.New(c),
	}
}

// RedisLocker implements the Locker interface for Redis
type RedisLocker struct {
	locker *redislock.Client
}

// RedisUnlocker implements the Unlocker interface for Redis
type RedisUnlocker struct {
	lock *redislock.Lock
}

// Unlock releases a Redis mutex
func (r *RedisUnlocker) Unlock() error {
	return r.lock.Release()
}

// Lock attempts to acquire a lock on a specific Redis key and sets expiry if acquired in case release is not triggered
func (l *RedisLocker) Lock(ctx context.Context, key string, opts ...Option) (Unlocker, error) {
	cfg := configWithOptions(opts)
	lock, err := l.locker.Obtain(key, cfg.ttl, &redislock.Options{
		Context:       ctx,
		RetryStrategy: redislock.NoRetry(),
	})
	if err != nil {
		return nil, err
	}
	return &RedisUnlocker{lock}, nil
}
