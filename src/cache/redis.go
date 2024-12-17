package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	c  Config
	cl *redis.Client
}

func NewRedisCache(ctx context.Context, c Config) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Password: c.Password,
	})
	_, err := rdb.Get(ctx, "").Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	return &RedisCache{c: c, cl: rdb}, nil
}

func (r *RedisCache) Set(ctx context.Context, key, url string) error {
	ctx, cancel := context.WithTimeout(ctx, r.c.WriteTimeout)
	defer cancel()
	return r.cl.Set(ctx, key, url, r.c.TTL).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.c.ReadTimeout)
	defer cancel()
	url, err := r.cl.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return url, nil
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, r.c.WriteTimeout)
	defer cancel()
	return r.cl.Del(ctx, key).Err()
}
