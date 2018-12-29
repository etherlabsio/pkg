package redis

import (
	"context"
	"time"

	"github.com/etherlabsio/errors"
	"github.com/etherlabsio/pkg/logutil"
	"github.com/go-kit/kit/log"

	"github.com/vmihailenco/msgpack"

	"github.com/go-redis/cache"
)

// Cache is an alternate implementation of a redis cache with built in pooling
type Cache struct {
	client    *Client
	codec     *cache.Codec
	logger    Logger
	namespace string
}

// NewCache returns a CacheV2 struct
func NewCache(client *Client, opts ...CacheOption) *Cache {
	option := NewCacheOptions(opts...)
	logger := log.With(client.logger, "component", "cache", "source", "redis")
	return &Cache{
		client: client,
		codec: &cache.Codec{
			Redis: client,
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)
			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		},
		namespace: option.Namespace,
		logger:    logger,
	}
}

// WithNamespace returns a new cache struct with a different namepace using the same underlying connection
func (c *Cache) WithNamespace(ns string) *Cache {
	cache := *c
	cache.namespace = ns
	return &cache
}

// Check is an implementation for our healthcheck endpoint to see if redis is responding to the ping
func (c *Cache) Check(_ context.Context) error {
	cmd := c.client.Ping()
	return cmd.Err()
}

// Get gets a value for a key along with the namespace from inside the cache
func (c *Cache) Get(k string, val interface{}) bool {
	const op errors.Op = "redis.Get"
	key := c.nsKey(k)
	err := c.codec.Get(key, val)
	return c.checkErr(op, key, err)
}

// Set sets a value for a key along with the namespace to the cache
func (c *Cache) Set(k string, val interface{}, exp time.Duration) bool {
	const op errors.Op = "redis.Set"
	key := c.nsKey(k)
	err := c.codec.Set(&cache.Item{
		Key:        key,
		Object:     val,
		Expiration: exp,
	})
	return c.checkErr(op, key, err)
}

// Delete deletes a value for a key along with the namespace from the cache
func (c *Cache) Delete(k string) bool {
	const op errors.Op = "redis.Set"
	key := c.nsKey(k)
	err := c.codec.Delete(key)
	return c.checkErr(op, key, err)
}

func (c *Cache) checkErr(op errors.Op, key string, err error) bool {
	if err == nil {
		return true
	}
	logutil.WithError(c.logger, err).Log("op", op, "key", key)
	return false
}

func (c *Cache) nsKey(k string) string {
	const separator = ":"
	return c.namespace + separator + k
}
