package service

import (
	"context"

	"lingua-evo/internal/db/redis"
	"lingua-evo/internal/services/session/entity"
)

type SessionSvc struct {
	redis *redis.Redis
}

func NewService(redis *redis.Redis) *SessionSvc {
	return &SessionSvc{
		redis: redis,
	}
}

func (s *SessionSvc) GetSession(ctx context.Context, sid string) (*entity.Session, error) {
	return nil, nil
}
