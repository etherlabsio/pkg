package locker

import (
	"context"
	"time"

	"github.com/bsm/redislock"
)

// DefaultExpiry signifies the default expiration for an acquired lock
const DefaultExpiry = 1 * time.Second

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = redislock.ErrNotObtained

	// ErrNotHeld is returned when trying to release an inactive lock.
	ErrNotHeld = redislock.ErrLockNotHeld
)

type Option func(*config)

func WithTTL(ttl time.Duration) Option {
	return func(cfg *config) {
		if ttl <= 0 {
			return
		}
		cfg.ttl = ttl
	}
}

// Locker creates a distributed lock provided a uniquely identifiable lock key
type Locker interface {
	Lock(context.Context, string, ...Option) (Unlocker, error)
}

// Unlocker releases the mutex lock
type Unlocker interface {
	Unlock() error
}

type config struct {
	ttl time.Duration
}

func configWithOptions(opts []Option) *config {
	cfg := &config{
		ttl: DefaultExpiry,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
