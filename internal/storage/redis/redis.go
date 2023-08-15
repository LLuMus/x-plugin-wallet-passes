package redis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(address, password, redisTLS string) *Cache {
	opt := &redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,
	}

	if redisTLS == "true" {
		opt.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(opt)
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return &Cache{
		rdb: rdb,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}
