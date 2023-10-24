package service

import (
	"context"
	"fmt"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/services/auth/dto"
	"lingua-evo/internal/services/auth/entity"
	entityUser "lingua-evo/internal/services/user/entity"
	"lingua-evo/pkg/jwt"
	"lingua-evo/pkg/tools"

	"github.com/google/uuid"
)

type (
	sessionRepo interface {
		SetSession(ctx context.Context, s *entity.Session) error
	}

	userSvc interface {
		GetUser(ctx context.Context, login string) (*entityUser.User, error)
	}

	AuthSvc struct {
		repo    sessionRepo
		userSvc userSvc
	}
)

func NewService(repo sessionRepo, userSvc userSvc) *AuthSvc {
	return &AuthSvc{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (s *AuthSvc) CreateSession(ctx context.Context, sessionRq *dto.CreateSessionRq) (*entity.Tokens, error) {
	u, err := s.userSvc.GetUser(ctx, sessionRq.User)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - getUser: %v", err)
	}
	if err := tools.CheckPasswordHash(sessionRq.Password, u.PasswordHash); err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - incorrect password: %v", err)
	}

	session := &entity.Session{
		ID:           uuid.New(),
		UserID:       u.ID,
		RefreshToken: uuid.New(),
		ExpiresAt:    time.Now().UTC().Add(time.Duration(config.GetConfig().JWT.ExpireAccess)),
		CreatedAt:    time.Now().UTC(),
	}

	err = s.repo.SetSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("auth.delivery.Handler.createSession - setSession: %v", err)
	}

	claims := &entity.Claims{
		ID:        session.ID,
		UserID:    u.ID,
		ExpiresAt: session.ExpiresAt,
	}

	jwtToken, err := jwt.NewJWTToken(u, claims)
	if err != nil {
		return nil, fmt.Errorf("auth.delivery.Handler.createSession - jwt.NewToken: %v", err)
	}

	tokens := &entity.Tokens{
		JWT:          jwtToken,
		RefreshToken: session.RefreshToken,
	}

	return tokens, nil
}
