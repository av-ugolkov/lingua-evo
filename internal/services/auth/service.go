package auth

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		SetAccountCode(ctx context.Context, email string, code int, expiration time.Duration) error
	}

	userSvc interface {
		GetUser(ctx context.Context, login string) (*entityUser.User, error)
		GetUserByEmail(ctx context.Context, email string) (*entityUser.User, error)
		UpdateLastVisited(ctx context.Context, uid uuid.UUID) error
	}

	emailSvc interface {
		SendAuthCode(toEmail string, code int) error
	}
)

type Service struct {
	repo    sessionRepo
	userSvc userSvc
	email   emailSvc
}

func NewService(repo sessionRepo, userSvc userSvc, email emailSvc) *Service {
	return &Service{
		repo:    repo,
		userSvc: userSvc,
		email:   email,
	}
}

func (s *Service) SignIn(ctx context.Context, user, password, fingerprint string, refreshTokenID uuid.UUID) (*Tokens, error) {
	u, err := s.userSvc.GetUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %w", err)
	}
	if err := utils.CheckPasswordHash(password, u.PasswordHash); err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn - [%w]: %v", ErrWrongPassword, err)
	}

	err = s.userSvc.UpdateLastVisited(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	now := time.Now().UTC()
	session := &Session{
		UserID:      u.ID,
		Fingerprint: fingerprint,
		CreatedAt:   now,
	}

	err = s.addRefreshSession(ctx, refreshTokenID, session)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	claims := &Claims{
		ID:        refreshTokenID,
		ExpiresAt: now.Add(duration),
	}

	accessToken, err := token.NewJWTToken(u.ID, claims.ID, claims.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	tokens := &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenID,
	}

	return tokens, nil
}

// RefreshSessionToken - the method is called from the client
func (s *Service) RefreshSessionToken(ctx context.Context, newTokenID, oldTokenID uuid.UUID, fingerprint string) (*Tokens, error) {
	oldRefreshSession, err := s.repo.GetSession(ctx, oldTokenID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	if oldRefreshSession.Fingerprint != fingerprint {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %w", errNotEqualFingerprints)
	}

	newSession := &Session{
		UserID:      oldRefreshSession.UserID,
		Fingerprint: oldRefreshSession.Fingerprint,
		CreatedAt:   time.Now().UTC(),
	}

	err = s.addRefreshSession(ctx, newTokenID, newSession)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	err = s.repo.DeleteSession(ctx, oldTokenID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	err = s.userSvc.UpdateLastVisited(ctx, oldRefreshSession.UserID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	claims := &Claims{
		ID:        newTokenID,
		ExpiresAt: time.Now().UTC().Add(duration),
	}

	accessToken, err := token.NewJWTToken(oldRefreshSession.UserID, claims.ID, claims.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession: %v", err)
	}

	tokens := &Tokens{
		AccessToken:  accessToken,
		RefreshToken: newTokenID,
	}

	return tokens, nil
}

func (s *Service) SignOut(ctx context.Context, refreshToken uuid.UUID, fingerprint string) error {
	oldRefreshSession, err := s.repo.GetSession(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("auth.Service.SignOut: %v", err)
	}

	if oldRefreshSession.Fingerprint != fingerprint {
		return fmt.Errorf("auth.Service.SignOut: %w", errNotEqualFingerprints)
	}

	err = s.repo.DeleteSession(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("auth.Service.SignOut: %v", err)
	}

	return nil
}

func (s *Service) CreateCode(ctx context.Context, email string) error {
	err := s.validateEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %v", err)
	}

	creatingCode := rand.IntN(999999-100000) + 100000

	err = s.email.SendAuthCode(email, creatingCode)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %v", err)
	}

	err = s.repo.SetAccountCode(ctx, email, creatingCode, time.Duration(10)*time.Minute)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %w", err)
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

func (s *Service) validateEmail(ctx context.Context, email string) error {
	if !utils.IsEmailValid(email) {
		return ErrEmailNotCorrect
	}

	userData, err := s.userSvc.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil
	} else if userData != nil && userData.ID == uuid.Nil && err == nil {
		return ErrItIsAdmin
	} else if userData != nil && userData.ID != uuid.Nil {
		return ErrEmailBusy
	}

	return nil
}
