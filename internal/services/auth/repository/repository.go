package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	authEntity "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	jsoniter "github.com/json-iterator/go"

	"github.com/google/uuid"
)

type (
	redis interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value any, expiration time.Duration) (string, error)
		SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error)
		GetTTL(ctx context.Context, key string) (time.Duration, error)
		Delete(ctx context.Context, key string) (int64, error)
	}

	SessionRepo struct {
		redis redis
	}
)

func NewRepo(r redis) *SessionRepo {
	return &SessionRepo{
		redis: r,
	}
}

func (r *SessionRepo) SetSession(ctx context.Context, key string, s *authEntity.Session, expiration time.Duration) error {
	data, err := s.JSON()
	if err != nil {
		return err
	}

	b, err := r.redis.SetNX(ctx, key, data, expiration)
	if !b {
		return err
	}
	return nil
}

func (r *SessionRepo) GetSession(ctx context.Context, key string) (*authEntity.Session, error) {
	value, err := r.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("auth.repository.SessionRepo.GetSession: %w", err)
	}

	var session authEntity.Session
	err = jsoniter.Unmarshal([]byte(value), &session)
	if err != nil {
		return nil, fmt.Errorf("auth.repository.SessionRepo.GetSession: %w", err)
	}
	return &session, nil
}

func (r *SessionRepo) GetCountSession(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	return count, nil
}

func (r *SessionRepo) DeleteSession(ctx context.Context, key string) error {
	_, err := r.redis.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("auth.repository.SessionRepo.DeleteSession: %w", err)
	}
	return nil
}

func (r *SessionRepo) DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (r *SessionRepo) SetAccountCode(ctx context.Context, email string, code int, ttl time.Duration) error {
	_, err := r.redis.SetNX(ctx, email, code, ttl)
	if err != nil {
		return fmt.Errorf("auth.repository.SessionRepo.SetAccountCode: %w", err)
	}

	return nil
}

func (r *SessionRepo) GetAccountCode(ctx context.Context, email string) (int, error) {
	value, err := r.redis.Get(ctx, email)
	if err != nil {
		return 0, fmt.Errorf("auth.repository.SessionRepo.GetAccountCode: %w", err)
	}

	code, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("auth.repository.SessionRepo.GetAccountCode: %w", err)
	}

	return code, nil
}
