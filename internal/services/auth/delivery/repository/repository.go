package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	authEntity "github.com/av-ugolkov/lingua-evo/internal/services/auth"

	"github.com/google/uuid"
)

type (
	redis interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value any, expiration time.Duration) (string, error)
		SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error)
		ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error)
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

func (r *SessionRepo) SetSession(ctx context.Context, tokenID uuid.UUID, s *authEntity.Session, expiration time.Duration) error {
	b, err := r.redis.SetNX(ctx, tokenID.String(), s, expiration)
	if !b {
		return err
	}
	return nil
}

func (r *SessionRepo) GetSession(ctx context.Context, refreshToken uuid.UUID) (*authEntity.Session, error) {
	var s authEntity.Session
	s2, err := r.redis.Get(ctx, refreshToken.String())
	if err != nil {
		return nil, fmt.Errorf("auth.repository.SessionRepo.GetSession: %w", err)
	}

	err = json.Unmarshal([]byte(s2), &s)
	if err != nil {
		return nil, fmt.Errorf("auth.repository.SessionRepo.GetSession - unmarshal: %w", err)
	}
	return &s, nil
}

func (r *SessionRepo) GetSessionExpire(ctx context.Context, uid, refreshTokenID uuid.UUID) (time.Time, error) {
	b, err := r.redis.ExpireAt(ctx, uid.String(), time.Now())
	if err != nil {
		return time.Now(), fmt.Errorf("auth.repository.SessionRepo.GetSessionExpire: %w", err)
	}

	//TODO check b
	fmt.Println(b)

	return time.Now(), nil
}

func (r *SessionRepo) GetCountSession(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	return count, nil
}

func (r *SessionRepo) DeleteSession(ctx context.Context, session uuid.UUID) error {
	_, err := r.redis.Delete(ctx, session.String())
	if err != nil {
		return fmt.Errorf("auth.repository.SessionRepo.DeleteSession: %w", err)
	}
	return nil
}

func (r *SessionRepo) DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (r *SessionRepo) SetAccountCode(ctx context.Context, email string, code int, expiration time.Duration) error {
	_, err := r.redis.Set(ctx, email, code, expiration)
	if err != nil {
		return fmt.Errorf("auth.repository.SessionRepo.SetAccountCode: %w", err)
	}

	return nil
}
