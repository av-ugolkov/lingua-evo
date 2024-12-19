package service

import (
	"context"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth/dto"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

func (s *Service) SignIn(ctx context.Context, user, password, fingerprint string, refreshTokenID uuid.UUID) (*dto.CreateSessionRs, string, error) {
	u, err := s.userSvc.GetUser(ctx, user)
	if err != nil {
		return nil, runtime.EmptyString, fmt.Errorf("auth.Service.SignIn: %w", err)
	}

	pswHash, err := s.userSvc.GetPswHash(ctx, u.ID)
	if err != nil {
		return nil, runtime.EmptyString, fmt.Errorf("auth.Service.SignIn: %w", err)
	}

	if err := utils.CheckPasswordHash(password, pswHash); err != nil {
		return nil, runtime.EmptyString, fmt.Errorf("auth.Service.SignIn - [%w]: %v", entity.ErrWrongPassword, err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	now := time.Now().UTC()
	session := &entity.Session{
		UserID:       u.ID,
		TypeToken:    entity.Email,
		RefreshToken: refreshTokenID.String(),
		ExpiresAt:    now.Add(duration),
	}

	err = s.addRefreshSession(ctx, fmt.Sprintf("%s:%s:%s", u.ID, fingerprint, RedisRefreshToken), session)
	if err != nil {
		return nil, runtime.EmptyString, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	accessToken, err := token.NewJWTToken(u.ID, refreshTokenID.String(), now.Add(duration))
	if err != nil {
		return nil, runtime.EmptyString, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	err = s.userSvc.UpdateVisitedAt(ctx, u.ID)
	if err != nil {
		return nil, runtime.EmptyString, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	tokens := &entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: session.RefreshToken,
	}

	return dto.CreateSessionToDTO(tokens), session.RefreshToken, nil
}

func (s *Service) SignUp(ctx context.Context, usr entity.User, fingerprint string) (uuid.UUID, error) {
	if err := s.validateEmail(ctx, usr.Email); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp: %v", err)
	}

	code, err := s.repo.GetAccountCode(ctx, fmt.Sprintf("%s:%s:%s", fingerprint, usr.Email, RedisCreateUser))
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp: %v", err)
	}

	if code != usr.Code {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp: code mismatch")
	}

	if err := s.validateUsername(ctx, usr.Nickname); err != nil {
		return uuid.Nil, err
	}

	if err := validatePassword(usr.Password); err != nil {
		return uuid.Nil, err
	}

	pswHash, err := utils.HashPassword(usr.Password)
	if err != nil {
		return uuid.Nil, err
	}

	uid, err := s.userSvc.AddUser(ctx, entityUser.User{
		Nickname: usr.Nickname,
		Email:    usr.Email,
		Role:     usr.Role,
	}, pswHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - AddUser: %v", err)
	}

	return uid, nil
}

func (s *Service) CreateCode(ctx context.Context, email string, fingerprint string) error {
	err := s.validateEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %v", err)
	}

	creatingCode := utils.GenerateCode()

	err = s.repo.SetAccountCode(ctx, fmt.Sprintf("%s:%s:%s", fingerprint, email, RedisCreateUser), creatingCode, time.Duration(5)*time.Minute)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %w", err)
	}

	err = s.email.SendAuthCode(email, creatingCode)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %v", err)
	}

	return nil
}

func (s *Service) refreshEmailToken(_ context.Context, uid uuid.UUID, refreshToken string) (*entity.Tokens, error) {
	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second

	accessToken, err := token.NewJWTToken(uid, refreshToken, time.Now().UTC().Add(duration))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession: %v", err)
	}

	tokens := &entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}
