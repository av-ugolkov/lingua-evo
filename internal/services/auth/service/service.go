package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/google"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	RedisCreateUser = "create_user"
)

type (
	sessionRepo interface {
		SetSession(ctx context.Context, key string, s entity.Session, ttl time.Duration) error
		GetSession(ctx context.Context, key string) (entity.Session, error)
		GetCountSession(ctx context.Context, userID uuid.UUID) (int64, error)
		DeleteSession(ctx context.Context, key string) error
		SetAccountCode(ctx context.Context, key string, code int, ttl time.Duration) error
		GetAccountCode(ctx context.Context, key string) (int, error)
	}

	userSvc interface {
		AddUser(ctx context.Context, userCreate entityUser.User, pswHash string) (uuid.UUID, error)
		AddGoogleUser(ctx context.Context, userCreate entityUser.GoogleUser) (uuid.UUID, error)
		GetUser(ctx context.Context, login string) (*entityUser.User, error)
		GetUserByNickname(ctx context.Context, nickname string) (*entityUser.User, error)
		GetPswHash(ctx context.Context, uid uuid.UUID) (string, error)
		GetUserByEmail(ctx context.Context, email string) (*entityUser.User, error)
		UpdateVisitedAt(ctx context.Context, uid uuid.UUID) error
		GetUserByGoogleID(ctx context.Context, googleID string) (*entityUser.User, error)
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

func (s *Service) RefreshSessionToken(ctx context.Context, newTokenID, oldTokenID any, fingerprint string) (*entity.Tokens, error) {
	oldRefreshSession, err := s.repo.GetSession(ctx, fmt.Sprintf("%s:%s", fingerprint, oldTokenID))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	var tokens *entity.Tokens
	switch oldTokenID.(type) {
	case string:
		accessToken, err := google.RefreshToken(ctx, oldTokenID.(string))
		if err != nil {
			return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
		}

		tokens = &entity.Tokens{
			AccessToken:  accessToken,
			RefreshToken: oldTokenID.(string),
		}
	case uuid.UUID:
		additionalTime := config.GetConfig().JWT.ExpireAccess
		duration := time.Duration(additionalTime) * time.Second
		accessToken, err := token.NewJWTToken(uuid.UUID(oldRefreshSession), newTokenID.(uuid.UUID), time.Now().UTC().Add(duration))
		if err != nil {
			return nil, fmt.Errorf("auth.Service.CreateSession: %v", err)
		}

		tokens = &entity.Tokens{
			AccessToken:  accessToken,
			RefreshToken: newTokenID.(uuid.UUID).String(),
		}
	default:
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", "unknown type")
	}

	err = s.addRefreshSession(ctx, fmt.Sprintf("%s:%s", fingerprint, tokens.RefreshToken), oldRefreshSession)
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

	return tokens, nil
}

func (s *Service) SignOut(ctx context.Context, refreshToken uuid.UUID, fingerprint string) error {
	err := s.repo.DeleteSession(ctx, fmt.Sprintf("%s:%s", fingerprint, refreshToken))
	if err != nil {
		return fmt.Errorf("auth.Service.SignOut: %v", err)
	}

	return nil
}

func (s *Service) GoogleAuthUrl() string {
	return google.GetAuthUrl()
}

func (s *Service) addRefreshSession(ctx context.Context, key string, refreshSession entity.Session) error {
	ttl := time.Duration(config.GetConfig().JWT.ExpireRefresh) * time.Second
	err := s.repo.SetSession(ctx, key, refreshSession, ttl)
	if err != nil {
		return fmt.Errorf("auth.Service.addRefreshSession: %w", err)
	}
	return nil
}

func (s *Service) validateEmail(ctx context.Context, email string) error {
	if !utils.IsEmailValid(email) {
		return entity.ErrEmailNotCorrect
	}

	userData, err := s.userSvc.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil
	} else if userData != nil && userData.ID == uuid.Nil && err == nil {
		return entity.ErrItIsAdmin
	} else if userData != nil && userData.ID != uuid.Nil {
		return entity.ErrEmailBusy
	}

	return nil
}

func (s *Service) validateUsername(ctx context.Context, username string) error {
	if len(username) <= entity.MinNicknameLen {
		return entity.ErrNicknameLen
	}

	userData, err := s.userSvc.GetUserByNickname(ctx, username)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil
	} else if userData.ID == uuid.Nil && err == nil {
		return entity.ErrItIsAdmin
	} else if userData.ID != uuid.Nil {
		return entity.ErrNicknameBusy
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < entity.MinPasswordLen {
		return entity.ErrPasswordLen
	}

	if !utils.IsPasswordValid(password) {
		return entity.ErrPasswordDifficult
	}

	return nil
}
