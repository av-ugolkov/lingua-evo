package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/pkg/token"
	"github.com/av-ugolkov/lingua-evo/pkg/utils"

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
		SetSession(ctx context.Context, tokenID uuid.UUID, s *Session, expiration time.Duration) error
		GetSession(ctx context.Context, refreshTokenID uuid.UUID) (*Session, error)
		GetCountSession(ctx context.Context, userID uuid.UUID) (int64, error)
		DeleteSession(ctx context.Context, session uuid.UUID) error
	}

	userSvc interface {
		GetUser(ctx context.Context, login string) (*entityUser.User, error)
		GetUserByID(ctx context.Context, uid uuid.UUID) (*entityUser.User, error)
	}

	Service struct {
		repo    sessionRepo
		userSvc userSvc
	}
)

func NewService(repo sessionRepo, userSvc userSvc) *Service {
	return &Service{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (s *Service) Login(ctx context.Context, user, password, fingerprint string) (*Tokens, error) {
	u, err := s.userSvc.GetUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession - getUser: %v", err)
	}
	if err := utils.CheckPasswordHash(password, u.PasswordHash); err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession - incorrect password: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	now := time.Now().UTC()
	session := &Session{
		UserID:      u.ID,
		Fingerprint: fingerprint,
		CreatedAt:   now,
	}

	refreshTokenID := uuid.New()
	err = s.addRefreshSession(ctx, refreshTokenID, session)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession - addRefreshSession: %v", err)
	}

	claims := &Claims{
		ID:        refreshTokenID,
		UserID:    u.ID,
		ExpiresAt: now.Add(duration),
	}

	accessToken, err := token.NewJWTToken(u.ID, claims.ID, claims.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession - jwt.NewToken: %v", err)
	}

	tokens := &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenID,
	}

	return tokens, nil
}

// RefreshSessionToken - the method is called from the client
func (s *Service) RefreshSessionToken(ctx context.Context, refreshToken uuid.UUID, fingerprint string) (*Tokens, error) {
	oldRefreshSession, err := s.repo.GetSession(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken - get session: %v", err)
	}

	if oldRefreshSession.Fingerprint != fingerprint {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %w", errNotEqualFingerprints)
	}

	tokenID := uuid.New()
	newSession := &Session{
		UserID:      oldRefreshSession.UserID,
		Fingerprint: oldRefreshSession.Fingerprint,
		CreatedAt:   time.Now().UTC(),
	}

	err = s.addRefreshSession(ctx, tokenID, newSession)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken - addRefreshSession: %v", err)
	}

	err = s.repo.DeleteSession(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken - delete session: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	claims := &Claims{
		ID:        tokenID,
		UserID:    oldRefreshSession.UserID,
		ExpiresAt: time.Now().UTC().Add(duration),
	}
	u, err := s.userSvc.GetUserByID(ctx, oldRefreshSession.UserID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession - get user by ID: %v", err)
	}

	accessToken, err := token.NewJWTToken(u.ID, claims.ID, claims.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession - create access token: %v", err)
	}

	tokens := &Tokens{
		AccessToken:  accessToken,
		RefreshToken: tokenID,
	}

	return tokens, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken uuid.UUID, fingerprint string) error {
	oldRefreshSession, err := s.repo.GetSession(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("auth.Service.Logout - GetSession: %v", err)
	}

	if oldRefreshSession.Fingerprint != fingerprint {
		return fmt.Errorf("auth.Service.Logout: %w", errNotEqualFingerprints)
	}

	err = s.repo.DeleteSession(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("auth.Service.Logout - DeleteSession: %v", err)
	}

	return nil
}

func (s *Service) addRefreshSession(ctx context.Context, tokenID uuid.UUID, refreshSession *Session) error {
	expiration := time.Duration(config.GetConfig().JWT.ExpireRefresh) * time.Second
	err := s.repo.SetSession(ctx, tokenID, refreshSession, expiration)
	if err != nil {
		return fmt.Errorf("auth.Service.addRefreshSession: %w", err)
	}
	return nil
}
