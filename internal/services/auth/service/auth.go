package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/services/auth/dto"
	"lingua-evo/internal/services/auth/entity"
	entityUser "lingua-evo/internal/services/user/entity"
	"lingua-evo/pkg/token"
	"lingua-evo/pkg/utils"

	"github.com/google/uuid"
)

const (
	MAX_REFRESH_SESSIONS_COUNT = 5
)

type (
	sessionRepo interface {
		SetSession(ctx context.Context, tokenID uuid.UUID, s *entity.Session, expiration time.Duration) error
		GetSession(ctx context.Context, userID, refreshTokenID uuid.UUID) (*entity.Session, error)
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

func (s *AuthSvc) Login(ctx context.Context, sessionRq *dto.CreateSessionRq) (*entity.Tokens, error) {
	u, err := s.userSvc.GetUser(ctx, sessionRq.User)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - getUser: %v", err)
	}
	if err := utils.CheckPasswordHash(sessionRq.Password, u.PasswordHash); err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - incorrect password: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	now := time.Now().UTC()
	session := &entity.Session{
		UserID:      u.ID,
		Fingerprint: sessionRq.Fingerprint,
		CreatedAt:   now,
	}

	refreshTokenID := uuid.New()
	err = s.addRefreshSession(ctx, refreshTokenID, session)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - addRefreshSession: %v", err)
	}

	claims := &entity.Claims{
		ID:        refreshTokenID,
		UserID:    u.ID,
		ExpiresAt: now.Add(duration),
	}

	accessToken, err := token.NewJWTToken(u, claims)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - jwt.NewToken: %v", err)
	}

	tokens := &entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenID,
	}

	return tokens, nil
}

// RefreshSessionToken - the method is called from the client
func (s *AuthSvc) RefreshSessionToken(ctx context.Context, uid, refreshToken uuid.UUID) (*entity.Tokens, error) {
	oldRefreshSession, err := s.repo.GetSession(ctx, uid, refreshToken)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - GetSession: %v", err)
	} else if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - GetSession: %v", err)
	}
	err = s.repo.DeleteSession(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - DeleteSession: %v", err)
	}

	err = s.verifyRefreshSession(ctx, uid, oldRefreshSession)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - verifyRefreshSession: %v", err)
	}

	tokenID := uuid.New()
	newSession := &entity.Session{
		UserID:      uid,
		Fingerprint: oldRefreshSession.Fingerprint,
		CreatedAt:   time.Now().UTC(),
	}

	err = s.addRefreshSession(ctx, tokenID, newSession)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - addRefreshSession: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	claims := &entity.Claims{
		ID:        tokenID,
		UserID:    uid,
		ExpiresAt: time.Now().UTC().Add(duration),
	}
	u, err := s.userSvc.GetUserByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - get user by ID: %v", err)
	}

	accessToken, err := token.NewJWTToken(u, claims)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - create access token: %v", err)
	}

	tokens := &entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: tokenID,
	}

	return tokens, nil
}

func (s *AuthSvc) Logout(ctx context.Context, uid uuid.UUID) error {
	err := s.repo.DeleteSession(ctx, uid)
	if err != nil {
		return fmt.Errorf("auth.service.AuthSvc.logout - DeleteSession: %v", err)
	}

	return nil
}

func (s *AuthSvc) addRefreshSession(ctx context.Context, tokenID uuid.UUID, refreshSession *entity.Session) error {
	expiration := time.Duration(config.GetConfig().JWT.ExpireRefresh) * time.Second
	err := s.repo.SetSession(ctx, tokenID, refreshSession, expiration)
	if err != nil {
		return fmt.Errorf("auth.service.AuthSvc.addRefreshSession: %w", err)
	}
	return nil
}

func (s *AuthSvc) verifyRefreshSession(ctx context.Context, uid uuid.UUID, oldRefreshSession *entity.Session) error {
	//TODO нужно получить все сессии пользователя и проверить сессию с определеным ID
	/*expireSession, err := s.repo.GetSessionExpire(ctx, uid, refreshTokenID)
	if err != nil {
		return fmt.Errorf("auth.service.AuthSvc.verifyRefreshSession - GetSession: %v", err)
	}
	if oldRefreshSession.ExpiresAt.Before(time.Now().UTC()) {
		return fmt.Errorf("auth.service.AuthSvc.verifyRefreshSession - session expired")
	}*/
	return nil
}
