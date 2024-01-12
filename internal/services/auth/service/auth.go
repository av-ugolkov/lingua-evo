package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"lingua-evo/internal/config"
	entity "lingua-evo/internal/services/auth"
	entityUser "lingua-evo/internal/services/user"
	"lingua-evo/pkg/token"
	"lingua-evo/pkg/utils"

	"github.com/google/uuid"
)

const (
	MAX_REFRESH_SESSIONS_COUNT = 5
)

var (
	errNotEqualFingerprints = errors.New("new fingerprint is not equal old fingerprint")
)

type (
	sessionRepo interface {
		SetSession(ctx context.Context, tokenID uuid.UUID, s *entity.Session, expiration time.Duration) error
		GetSession(ctx context.Context, refreshTokenID uuid.UUID) (*entity.Session, error)
		GetCountSession(ctx context.Context, userID uuid.UUID) (int64, error)
		DeleteSession(ctx context.Context, session uuid.UUID) error
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

func (s *AuthSvc) Login(ctx context.Context, user, password, fingerprint string) (*entity.Tokens, error) {
	u, err := s.userSvc.GetUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - getUser: %v", err)
	}
	if err := utils.CheckPasswordHash(password, u.PasswordHash); err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.CreateSession - incorrect password: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	now := time.Now().UTC()
	session := &entity.Session{
		UserID:      u.ID,
		Fingerprint: fingerprint,
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
func (s *AuthSvc) RefreshSessionToken(ctx context.Context, refreshToken uuid.UUID, fingerprint string) (*entity.Tokens, error) {
	oldRefreshSession, err := s.repo.GetSession(ctx, refreshToken)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - GetSession: %v", err)
	} else if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - GetSession: %v", err)
	}

	if oldRefreshSession.Fingerprint != fingerprint {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken: %w", errNotEqualFingerprints)
	}

	tokenID := uuid.New()
	newSession := &entity.Session{
		UserID:      oldRefreshSession.UserID,
		Fingerprint: oldRefreshSession.Fingerprint,
		CreatedAt:   time.Now().UTC(),
	}

	err = s.addRefreshSession(ctx, tokenID, newSession)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - addRefreshSession: %v", err)
	}

	err = s.repo.DeleteSession(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth.service.AuthSvc.RefreshSessionToken - delete session: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	claims := &entity.Claims{
		ID:        tokenID,
		UserID:    oldRefreshSession.UserID,
		ExpiresAt: time.Now().UTC().Add(duration),
	}
	u, err := s.userSvc.GetUserByID(ctx, oldRefreshSession.UserID)
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
