package rediscache

import (
	"context"
	"errors"
	"log"
	"math/rand/v2"
	"time"

	"github.com/redis/go-redis/v9"
)

type options struct {
	ttl    time.Duration
	jitter time.Duration
}

type Option func(*options)

func WithExpiration(t time.Duration) Option {
	return func(o *options) {
		o.ttl = t
	}
}

func WithJitter(t time.Duration) Option {
	return func(o *options) {
		o.jitter = t
	}
}

type RedisCache struct {
	rds    redis.UniversalClient
	ttl    time.Duration
	jitter time.Duration
}

func NewRedisCache(client redis.UniversalClient, opts ...Option) *RedisCache {
	opt := &options{
		ttl:    0, // no exp on default
		jitter: 0, // no jitter on default
	}

	for _, optFn := range opts {
		optFn(opt)
	}

	return &RedisCache{
		rds:    client,
		ttl:    opt.ttl,
		jitter: opt.jitter,
	}
}

func (r *RedisCache) Get(key string) (data []byte, ok bool) {
	b, err := r.rds.Get(context.TODO(), key).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Default().Printf("rediscache: error when getting key %s: %s", key, err.Error())
		}
		return
	}

	return b, true
}

func (r *RedisCache) Set(key string, data []byte) {
	ttl := r.ttl
	if ttl > 0 && r.jitter > 0 {
		ttl += time.Duration(rand.Int64N(r.jitter.Nanoseconds()))
	}

	err := r.rds.Set(context.TODO(), key, data, ttl).Err()
	if err != nil {
		log.Default().Printf("rediscache: error when setting key %s: %s", key, err.Error())
	}
}

func (r *RedisCache) Delete(key string) {
	err := r.rds.Del(context.TODO(), key).Err()
	if err != nil {
		log.Default().Printf("rediscache: error when deleting key %s: %s", key, err.Error())
	}
}
