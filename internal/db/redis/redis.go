package redis

import (
	"fmt"

	"lingua-evo/internal/config"

	rClient "github.com/redis/go-redis/v9"
)

type Redis struct {
	r *rClient.Client
}

func New(cfg *config.Config) *Redis {
	r := rClient.NewClient(&rClient.Options{
		Addr:       fmt.Sprintf("%s:%s", cfg.DbRedis.Host, cfg.DbRedis.Port),
		ClientName: cfg.DbRedis.Name,
		Username:   cfg.DbRedis.User,
		Password:   cfg.DbRedis.Password,
		DB:         0,
	})
	return &Redis{
		r: r,
	}
}
