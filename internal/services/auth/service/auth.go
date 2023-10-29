package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/services/auth/dto"
	"lingua-evo/internal/services/auth/entity"
	entityUser "lingua-evo/internal/services/user/entity"
	"lingua-evo/pkg/jwt"
	"lingua-evo/pkg/tools"

	"github.com/google/uuid"
)

const (
	MAX_REFRESH_SESSIONS_COUNT = 5
)

type (
	sessionRepo interface {
		SetSession(ctx context.Context, s *entity.Session) error
		GetSession(ctx context.Context, refreshToken uuid.UUID) (*entity.Session, error)
		GetCountSession(ctx context.Context, userID uuid.UUID) (int64, error)
		DeleteSession(ctx context.Context, session uuid.UUID) error
		DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error
	}

	userSvc interface {
		GetUser(ctx context.Context, login string) (*entityUser.User, error)
		GetUserByID(ctx context.Context, uid uuid.UUID) (*entityUser.User, error)
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

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	session := &entity.Session{
		RefreshToken: uuid.New(),
		ExpiresAt:    time.Now().UTC().Add(duration),
		CreatedAt:    time.Now().UTC(),
		UserID:       u.ID,
	}

	if s.validSessionCount(ctx, session.UserID) {
		err = s.addRefreshSession(ctx, session)
		if err != nil {
			return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - addRefreshSession: %v", err)
		}
	} else {
		err := s.wipeAllUserRefreshSessions(ctx, session.UserID)
		if err != nil {
			return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - wipeAllUserRefreshSessions: %v", err)
		}
		err = s.addRefreshSession(ctx, session)
		if err != nil {
			return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - addRefreshSession after wipe: %v", err)
		}
	}

	claims := &entity.Claims{
		ID:        session.RefreshToken,
		UserID:    u.ID,
		ExpiresAt: session.ExpiresAt,
	}

	accessToken, err := jwt.NewJWTToken(u, claims)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - jwt.NewToken: %v", err)
	}

	tokens := &entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: session.RefreshToken,
	}

	return tokens, nil
}

func (s *AuthSvc) RefreshSessionToken(ctx context.Context, refreshToken uuid.UUID) (*entity.Tokens, error) {
	oldRefreshSession, err := s.repo.GetSession(ctx, refreshToken)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - GetSession: %v", err)
	} else if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - GetSession: %v", err)
	}
	err = s.repo.DeleteSession(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - DeleteSession: %v", err)
	}

	err = s.verifyRefreshSession(oldRefreshSession)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - verifyRefreshSession: %v", err)
	}

	duration := time.Second * time.Duration(config.GetConfig().JWT.ExpireRefresh)
	newSession := &entity.Session{
		RefreshToken: uuid.New(),
		ExpiresAt:    time.Now().UTC().Add(duration),
		CreatedAt:    time.Now().UTC(),
		UserID:       oldRefreshSession.UserID,
	}

	err = s.addRefreshSession(ctx, newSession)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - addRefreshSession: %v", err)
	}

	claims := &entity.Claims{
		ID:        newSession.RefreshToken,
		UserID:    newSession.UserID,
		ExpiresAt: newSession.ExpiresAt,
	}
	u, err := s.userSvc.GetUserByID(ctx, newSession.UserID)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - get user by ID: %v", err)
	}

	accessToken, err := jwt.NewJWTToken(u, claims)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - create access token: %v", err)
	}

	tokens := &entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: newSession.RefreshToken,
	}

	return tokens, nil
}

func (s *AuthSvc) validSessionCount(ctx context.Context, uid uuid.UUID) bool {
	i, err := s.repo.GetCountSession(ctx, uid)
	if err != nil {
		slog.Warn(fmt.Sprintf("auth.delivery.Handler.createSession - GetCountSession: %v", err))
		return false
	}
	return i < MAX_REFRESH_SESSIONS_COUNT
}

func (s *AuthSvc) addRefreshSession(ctx context.Context, refreshSession *entity.Session) error {
	err := s.repo.SetSession(ctx, refreshSession)
	if err != nil {
		return fmt.Errorf("auth.service.AuthSvc.addRefreshSession: %w", err)
	}
	return nil
}

func (s *AuthSvc) wipeAllUserRefreshSessions(ctx context.Context, uid uuid.UUID) error {
	err := s.repo.DeleteAllUserSessions(ctx, uid)
	if err != nil {
		return fmt.Errorf("auth.service.AuthSvc.wipeAllUserRefreshSessions: %w", err)
	}
	return nil
}

func (s *AuthSvc) verifyRefreshSession(oldRefreshSession *entity.Session) error {
	if oldRefreshSession.ExpiresAt.Before(time.Now().UTC()) {
		return fmt.Errorf("auth.service.AuthSvc.verifyRefreshSession - session expired")
	}
	return nil
}

/*
login
  get user
  check password
  create refresh session
  add refresh session
    if valid session
      add refresh session
    else
      wipe all refresh session
      add refresh session
  make access token use id, name, role, email
  create data from acces token and refresh token


const newRefreshSession = new RefreshSessionEntity({
      refreshToken: uuidv4(),
      userId: user.id,
      ip: ctx.ip,
      ua: ctx.headers['User-Agent'],
      fingerprint: ctx.body.fingerprint,
      expiresIn: refTokenExpiresInMilliseconds
    })

return this.result({
      data: {
        accessToken: await makeAccessToken(user),
        refreshToken: newRefreshSession.refreshToken
      },
      cookies: [
        new CookieEntity({
          name: 'refreshToken',
          value: newRefreshSession.refreshToken,
          domain: 'localhost',
          path: '/auth',
          maxAge: refTokenExpiresInSeconds,
          secure: false // temp: should be deleted
        })
      ]
    })
*/
