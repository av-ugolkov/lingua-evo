package auth

import (
	"context"
	"errors"
	"fmt"
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
		SetSession(ctx context.Context, key string, s Session, ttl time.Duration) error
		GetSession(ctx context.Context, key string) (Session, error)
		GetCountSession(ctx context.Context, userID uuid.UUID) (int64, error)
		DeleteSession(ctx context.Context, key string) error
		SetAccountCode(ctx context.Context, email string, code int, ttl time.Duration) error
	}

	userSvc interface {
		GetUser(ctx context.Context, login string) (*entityUser.User, error)
		GetPswHash(ctx context.Context, uid uuid.UUID) (string, error)
		GetUserByEmail(ctx context.Context, email string) (*entityUser.User, error)
		UpdateVisitedAt(ctx context.Context, uid uuid.UUID) error
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

	pswHash, err := s.userSvc.GetPswHash(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %w", err)
	}

	if err := utils.CheckPasswordHash(password, pswHash); err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn - [%w]: %v", ErrWrongPassword, err)
	}

	err = s.userSvc.UpdateVisitedAt(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	now := time.Now().UTC()
	session := Session(u.ID)

	err = s.addRefreshSession(ctx, fmt.Sprintf("%s:%s", fingerprint, refreshTokenID), session)
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
	oldRefreshSession, err := s.repo.GetSession(ctx, fmt.Sprintf("%s:%s", fingerprint, oldTokenID))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	err = s.addRefreshSession(ctx, fmt.Sprintf("%s:%s", fingerprint, newTokenID), oldRefreshSession)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	err = s.repo.DeleteSession(ctx, fmt.Sprintf("%s:%s", fingerprint, oldTokenID))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	err = s.userSvc.UpdateVisitedAt(ctx, uuid.UUID(oldRefreshSession))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	claims := &Claims{
		ID:        newTokenID,
		ExpiresAt: time.Now().UTC().Add(duration),
	}

	accessToken, err := token.NewJWTToken(uuid.UUID(oldRefreshSession), claims.ID, claims.ExpiresAt)
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
	err := s.repo.DeleteSession(ctx, fmt.Sprintf("%s:%s", fingerprint, refreshToken))
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

	creatingCode := utils.GenerateCode()

	err = s.email.SendAuthCode(email, creatingCode)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %v", err)
	}

	err = s.repo.SetAccountCode(ctx, email, creatingCode, time.Duration(5)*time.Minute)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %w", err)
	}

	return nil
}

func (s *Service) addRefreshSession(ctx context.Context, key string, refreshSession Session) error {
	ttl := time.Duration(config.GetConfig().JWT.ExpireRefresh) * time.Second
	err := s.repo.SetSession(ctx, key, refreshSession, ttl)
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
