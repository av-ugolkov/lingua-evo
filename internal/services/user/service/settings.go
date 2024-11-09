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
)

func (s *Service) SendSecurityCodeForUpdatePsw(ctx context.Context, uid uuid.UUID, psw, typeCode string) error {
	usr, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), ErrMsgUserNotFound)
	}

	if utils.CheckPasswordHash(psw, usr.PasswordHash) != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: incorrect password"), ErrMsgIncorrectPsw)
	}

	code := utils.GenerateCode()

	value, err := s.redis.SetNX(ctx, fmt.Sprintf("%s:%s", uid, typeCode), code, 5*time.Minute)
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
