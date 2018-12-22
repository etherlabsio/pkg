package redis

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis"
)

type Logger interface {
	Log(keyvals ...interface{}) error
}

// Options defines the set of parameters that can be passed as optional
type Options struct {
	Addresses []string
	Password  string
}

type CacheOptions struct {
	Logger    Logger
	Namespace string
}

type Option func(*Options)

type CacheOption func(*CacheOptions)

func Addresses(addrs ...string) Option {
	return func(opt *Options) {
		opt.Addresses = addrs
	}
}

func Password(pass string) Option {
	return func(opt *Options) {
		opt.Password = pass
	}
}

func Namespace(ns string) CacheOption {
	return func(opt *CacheOptions) {
		opt.Namespace = ns
	}
}

// NewOptions returns an Options struct with default options set
func NewOptions(opts ...Option) Options {
	option := Options{
		Addresses: []string{":6379"},
		Password:  "",
	}

	for _, opt := range opts {
		opt(&option)
	}
	return option
}

func NewCacheOptions(opts ...CacheOption) CacheOptions {
	option := CacheOptions{
		Logger:    log.NewNopLogger(),
		Namespace: "default",
	}

	for _, opt := range opts {
		opt(&option)
	}
	return option
}

type Client redis.UniversalClient

func NewClient(opts ...Option) Client {
	option := NewOptions(opts...)
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    option.Addresses,
		Password: option.Password,
		OnConnect: func(conn *redis.Conn) error {
			fmt.Println("connected to redis")
			return nil
		},
		PoolSize:     10,
		MaxRetries:   2,
		MinIdleConns: 5,
	})
}
