package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/runtime"

	rClient "github.com/redis/go-redis/v9"
)

type Redis struct {
	client *rClient.Client
}

func New(cfg *config.Config) *Redis {
	r := rClient.NewClient(&rClient.Options{
		Addr:       fmt.Sprintf("%s:%s", cfg.DbRedis.Host, cfg.DbRedis.Port),
		ClientName: cfg.DbRedis.Name,
		Password:   cfg.DbRedis.Password,
		DB:         0,
	})
	return &Redis{
		client: r,
	}
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return runtime.EmptyString, fmt.Errorf("redis.Get - key [%s]: %w", key, err)
	}

	return value, nil
}

func (r *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) (string, error) {
	return r.client.Set(ctx, key, value, expiration).Result()
}

func (r *Redis) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

func (r *Redis) Delete(ctx context.Context, key string) (int64, error) {
	return r.client.Del(ctx, key).Result()
}

func (r *Redis) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}
