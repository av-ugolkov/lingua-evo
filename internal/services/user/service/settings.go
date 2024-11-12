package service

import (
	"context"
	"fmt"
	"time"

	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"

	"github.com/google/uuid"
)

const (
	ErrMsgUserNotFound   = "Sorry, user not found"
	ErrMsgIncorrectPsw   = "Incorrect password"
	ErrMsgSamePsw        = "The same password"
	ErrMsgIncorrectEmail = "Incorrect email"
	ErrMsgSameEmail      = "The same email"
	ErrMsgInvalidEmail   = "Invalid email"
)

const (
	RedisUpdatePsw   = "update_psw"
	RedisUpdateEmail = "update_email"
)

type (
	repoSettings interface {
		GetPswHash(ctx context.Context, uid uuid.UUID) (string, error)
		UpdatePsw(ctx context.Context, uid uuid.UUID, hashPsw string) (err error)
		UpdateEmail(ctx context.Context, uid uuid.UUID, newEmail string) (err error)
	}

	emailSvc interface {
		SendEmailForUpdatePassword(toEmail, userName string, code int) error
		SendEmailForUpdateEmail(toEmail, userName string, code int) error
	}
)

func (s *Service) GetPswHash(ctx context.Context, uid uuid.UUID) (string, error) {
	pswHash, err := s.repo.GetPswHash(ctx, uid)
	if err != nil {
		return pswHash, fmt.Errorf("auth.Service.GetPswHash: %v", err)
	}

	return pswHash, nil
}

func (s *Service) SendSecurityCodeForUpdatePsw(ctx context.Context, uid uuid.UUID, psw string) error {
	pswHash, err := s.repo.GetPswHash(ctx, uid)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), ErrMsgUserNotFound)
	}

	if utils.CheckPasswordHash(psw, pswHash) != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: incorrect password"), ErrMsgIncorrectPsw)
	}

	code := utils.GenerateCode()
	value, err := s.redis.SetNX(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdatePsw), code, 5*time.Minute)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), msgerr.ErrMsgInternal)
	}
	if !value {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), "You have already sent a code. Please wait.")
	}

	usr, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), ErrMsgUserNotFound)
	}

	err = s.emailSvc.SendEmailForUpdatePassword(usr.Email, usr.Nickname, code)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), msgerr.ErrMsgInternal)
	}

	return nil
}

func (s *Service) UpdatePsw(ctx context.Context, uid uuid.UUID, oldPsw, newPsw, code string) error {
	pswHash, err := s.repo.GetPswHash(ctx, uid)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), ErrMsgUserNotFound)
	}

	if utils.CheckPasswordHash(oldPsw, pswHash) != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: incorrect password"), ErrMsgIncorrectPsw)
	}

	redisCode, err := s.redis.Get(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdatePsw))
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: %w", err), msgerr.ErrMsgInternal)
	}
	if redisCode != code {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: incorrect code"), "Incorrect code")
	}

	isValid := utils.IsPasswordValid(newPsw)
	if !isValid {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: invalid password"), "Password is invalid. Password must be at least 8-20 characters long and contain at least one uppercase letter, one lowercase letter, and one number.")
	}

	if utils.CheckPasswordHash(newPsw, pswHash) == nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: the same password"), ErrMsgSamePsw)
	}

	hashPassword, err := utils.HashPassword(newPsw)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: %w", err), msgerr.ErrMsgInternal)
	}

	err = s.repo.UpdatePsw(ctx, uid, hashPassword)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: %w", err), msgerr.ErrMsgInternal)
	}

	//TODO: нужно скидывать сессию

	return nil
}

func (s *Service) SendSecurityCodeForUpdateEmail(ctx context.Context, uid uuid.UUID) error {
	usr, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdateEmail: %w", err), ErrMsgUserNotFound)
	}

	code := utils.GenerateCode()
	value, err := s.redis.SetNX(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdateEmail), code, 5*time.Minute)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdateEmail: %w", err), msgerr.ErrMsgInternal)
	}
	if !value {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdateEmail: %w", err), "You have already sent a code. Please wait.")
	}

	err = s.emailSvc.SendEmailForUpdateEmail(usr.Email, usr.Nickname, code)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdateEmail: %w", err), msgerr.ErrMsgInternal)
	}

	return nil
}

func (s *Service) UpdateEmail(ctx context.Context, uid uuid.UUID, newEmail, code string) error {
	usr, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: %w", err), ErrMsgUserNotFound)
	}

	if newEmail == usr.Email {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: the same email"), ErrMsgSameEmail)
	}

	if !utils.IsEmailValid(newEmail) {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: invalid email"), ErrMsgInvalidEmail)
	}

	redisCode, err := s.redis.Get(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdateEmail))
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: %w", err), msgerr.ErrMsgInternal)
	}
	if redisCode != code {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: incorrect code"), "Incorrect code")
	}

	err = s.repo.UpdateEmail(ctx, uid, newEmail)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: %w", err), msgerr.ErrMsgInternal)
	}

	return nil
}
