package service

import (
	"context"
	"encoding/json"
	"fmt"

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
	data, err := s.redis.Get(ctx, sid)
	if err != nil {
		return nil, fmt.Errorf("session.service.SessionSvc.GetSession: %w", err)
	}

	var session entity.Session
	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return nil, fmt.Errorf("session.service.SessionSvc.GetSession - umarshal: %w", err)
	}
	return &session, nil
}
