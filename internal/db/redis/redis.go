package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"

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
		return "", fmt.Errorf("redis.Get - key [%s]: %w", key, err)
	}
	return value, nil
}

func (r *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) (string, error) {
	return r.client.Set(ctx, key, value, expiration).Result()
}

func (r *Redis) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

func (r *Redis) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return r.client.ExpireAt(ctx, key, tm).Result()
}

func (r *Redis) Delete(ctx context.Context, key string) (int64, error) {
	return r.client.Del(ctx, key).Result()
}

func (r *Redis) GetAccountCode(ctx context.Context, email string) (int, error) {
	codeStr, err := r.Get(ctx, email)
	if err != nil {
		return 0, fmt.Errorf("redis.GetAccountCode: %w", err)
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return 0, fmt.Errorf("redis.GetAccountCode - convert string to int: %w", err)
	}

	return code, nil
}
