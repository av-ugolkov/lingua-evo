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
	ErrMsgUserNotFound = "Sorry,user not found"
	ErrMsgIncorrectPsw = "Incorrect password"
	ErrMsgSamePsw      = "The same password"
)

const (
	RedisUpdatePsw = "update_psw"
)

type (
	repoSettings interface {
		UpdatePsw(ctx context.Context, uid uuid.UUID, hashPsw string) (err error)
	}
)

func (s *Service) SendSecurityCodeForUpdatePsw(ctx context.Context, uid uuid.UUID, psw string) error {
	usr, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), ErrMsgUserNotFound)
	}

	if utils.CheckPasswordHash(psw, usr.PasswordHash) != nil {
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

	err = s.emailSvc.SendEmailForUpdatePassword(usr.Email, usr.Name, code)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), msgerr.ErrMsgInternal)
	}

	return nil
}

func (s *Service) UpdatePsw(ctx context.Context, uid uuid.UUID, oldPsw, newPsw, code string) error {
	usr, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: %w", err), ErrMsgUserNotFound)
	}

	if utils.CheckPasswordHash(oldPsw, usr.PasswordHash) != nil {
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

	if utils.CheckPasswordHash(newPsw, usr.PasswordHash) == nil {
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

	//TODO: нужно ли скидывать сессию?

	return nil
}
