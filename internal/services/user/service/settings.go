package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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
		UpdateNickname(ctx context.Context, uid uuid.UUID, newNickname string) (err error)
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

func (s *Service) SendSecurityCodeForUpdatePsw(ctx context.Context, uid uuid.UUID, psw string) (int, error) {
	dur := time.Duration(5 * time.Minute)
	pswHash, err := s.repo.GetPswHash(ctx, uid)
	if err != nil {
		return int(dur.Milliseconds()), msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), entity.ErrMsgUserNotFound)
	}

	if utils.CheckPasswordHash(psw, pswHash) != nil {
		return int(dur.Milliseconds()), msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: incorrect password"), entity.ErrMsgIncorrectPsw)
	}

	code := utils.GenerateCode()
	value, err := s.redis.SetNX(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdatePsw), code, 5*time.Minute)
	if err != nil {
		return int(dur.Milliseconds()), msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), msgerr.ErrMsgInternal)
	}
	if !value {
		dur, _ = s.redis.GetTTL(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdatePsw))
		return int(dur.Milliseconds()),
			msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", entity.ErrDuplicateCode),
				fmt.Sprintf(entity.ErrMsgDuplicateCode, dur.String()))
	}

	usr, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return int(dur.Milliseconds()), msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), entity.ErrMsgUserNotFound)
	}

	err = s.emailSvc.SendEmailForUpdatePassword(usr.Email, usr.Nickname, code)
	if err != nil {
		return int(dur.Milliseconds()), msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), msgerr.ErrMsgInternal)
	}

	return int(dur.Milliseconds()), nil
}

func (s *Service) UpdatePsw(ctx context.Context, uid uuid.UUID, oldPsw, newPsw, code string) error {
	pswHash, err := s.repo.GetPswHash(ctx, uid)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdatePsw: %w", err), entity.ErrMsgUserNotFound)
	}

	if utils.CheckPasswordHash(oldPsw, pswHash) != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: incorrect password"), entity.ErrMsgIncorrectPsw)
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
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: the same password"), entity.ErrMsgSamePsw)
	}

	hashPassword, err := utils.HashPassword(newPsw)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: %w", err), msgerr.ErrMsgInternal)
	}

	err = s.repo.UpdatePsw(ctx, uid, hashPassword)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdatePsw: %w", err), msgerr.ErrMsgInternal)
	}

	_, _ = s.redis.Delete(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdatePsw))

	//TODO: нужно скидывать сессию

	return nil
}

func (s *Service) SendSecurityCodeForUpdateEmail(ctx context.Context, uid uuid.UUID) (int, error) {
	dur := time.Duration(5 * time.Minute)
	usr, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return int(dur.Milliseconds()), msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdateEmail: %w", err), entity.ErrMsgUserNotFound)
	}

	code := utils.GenerateCode()
	value, err := s.redis.SetNX(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdateEmail), code, 5*time.Minute)
	if err != nil {
		return int(dur.Milliseconds()), msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdateEmail: %w", err), msgerr.ErrMsgInternal)
	}
	if !value {
		dur, _ = s.redis.GetTTL(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdateEmail))
		return int(dur.Milliseconds()),
			msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdateEmail: %w", entity.ErrDuplicateCode),
				fmt.Sprintf(entity.ErrMsgDuplicateCode, dur.String()))
	}

	err = s.emailSvc.SendEmailForUpdateEmail(usr.Email, usr.Nickname, code)
	if err != nil {
		return int(dur.Milliseconds()), msgerr.New(fmt.Errorf("auth.Service.SendSecurityCodeForUpdateEmail: %w", err), msgerr.ErrMsgInternal)
	}

	return int(dur.Milliseconds()), nil
}

func (s *Service) UpdateEmail(ctx context.Context, uid uuid.UUID, newEmail, code string) error {
	usr, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: %w", err), entity.ErrMsgUserNotFound)
	}

	if newEmail == usr.Email {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: the same email"), entity.ErrMsgSameEmail)
	}

	if !utils.IsEmailValid(newEmail) {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: invalid email"), entity.ErrMsgInvalidEmail)
	}

	redisCode, err := s.redis.Get(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdateEmail))
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: %w", err), msgerr.ErrMsgInternal)
	}
	if redisCode != code {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: incorrect code"), "Incorrect code")
	}

	err = s.repo.UpdateEmail(ctx, uid, newEmail)
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: %w", err), entity.ErrMsgBusyEmail)
	case err != nil:
		return msgerr.New(fmt.Errorf("auth.Service.UpdateEmail: %w", err), msgerr.ErrMsgInternal)
	}

	_, _ = s.redis.Delete(ctx, fmt.Sprintf("%s:%s", uid, RedisUpdatePsw))

	return nil
}

func (s *Service) UpdateNickname(ctx context.Context, uid uuid.UUID, newNickname string) error {
	if !utils.IsNicknameValid(newNickname) {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateNickname: invalid nickname"), entity.ErrMsgInvalidNickname)
	}

	tempNickname := strings.ToLower(newNickname)
	if strings.Contains(tempNickname, "admin") || strings.Contains(tempNickname, "moderator") {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateNickname: contains admin"), entity.ErrFobiddenNickname)
	}

	err := s.repo.UpdateNickname(ctx, uid, newNickname)
	if err != nil {
		return msgerr.New(fmt.Errorf("auth.Service.UpdateNickname: %w", err), msgerr.ErrMsgInternal)
	}

	return nil
}
