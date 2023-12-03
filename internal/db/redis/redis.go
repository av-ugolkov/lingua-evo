package redis

import (
	"context"
	"fmt"
	"time"

	"lingua-evo/internal/config"

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
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

func (r *Redis) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return r.client.ExpireAt(ctx, key, tm).Result()
}
