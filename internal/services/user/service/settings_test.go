package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	uid     = uuid.Nil
	entUser = &entity.User{
		Email:        "test@test.com",
		PasswordHash: "$2a$11$0CdrPvkPtvA2b1KGilY1l.E7XKELBrGL7i0obCFR9h2kW0BgCGgeu",
	}
)

func TestService_SendSecurityCodeForUpdatePsw(t *testing.T) {
	ctx := context.Background()
	psw := "psw"

	t.Run("not found user", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(nil, fmt.Errorf("error"))

		service := NewService(nil, repo, nil, nil)
		err := service.SendSecurityCodeForUpdatePsw(ctx, uid, psw)
		assert.Error(t, err)
	})

	t.Run("different psw", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		service := NewService(nil, repo, nil, nil)
		wrongPsw := "wrong_psw"
		err := service.SendSecurityCodeForUpdatePsw(ctx, uid, wrongPsw)
		assert.Error(t, err)
	})

	t.Run("set code in redis", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("SetNX", ctx, mock.Anything, mock.Anything, 5*time.Minute).Return(false, fmt.Errorf("error"))
		service := NewService(nil, repo, redisDB, nil)
		err := service.SendSecurityCodeForUpdatePsw(ctx, uid, psw)
		assert.Error(t, err)
	})

	t.Run("code already exist in redis", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("SetNX", ctx, mock.Anything, mock.Anything, 5*time.Minute).Return(false, nil)
		service := NewService(nil, repo, redisDB, nil)
		err := service.SendSecurityCodeForUpdatePsw(ctx, uid, psw)
		assert.Error(t, err)
	})

	t.Run("sending code is failed", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("SetNX", ctx, mock.Anything, mock.Anything, 5*time.Minute).Return(true, nil)
		emailMockSvc := new(mockEmailSvc)
		emailMockSvc.On("SendEmailForUpdatePassword", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
		service := NewService(nil, repo, redisDB, emailMockSvc)
		err := service.SendSecurityCodeForUpdatePsw(ctx, uid, psw)
		assert.Error(t, err)
	})

	t.Run("sending code is done", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("SetNX", ctx, mock.Anything, mock.Anything, 5*time.Minute).Return(true, nil)
		emailMockSvc := new(mockEmailSvc)
		emailMockSvc.On("SendEmailForUpdatePassword", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		service := NewService(nil, repo, redisDB, emailMockSvc)
		err := service.SendSecurityCodeForUpdatePsw(ctx, uid, psw)
		assert.NoError(t, err)
	})
}

func TestService_UpdatePsw(t *testing.T) {
	ctx := context.Background()
	oldPsw := "psw"
	newPsw := "new_psw_12"
	code := "1234"

	t.Run("not found redis code", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("Get", ctx, mock.Anything).Return("", fmt.Errorf("error"))
		service := NewService(nil, repo, redisDB, nil)
		err := service.UpdatePsw(ctx, uid, oldPsw, newPsw, code)
		assert.Error(t, err)
	})

	t.Run("redis code is different", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("Get", ctx, mock.Anything).Return("1233", nil)
		service := NewService(nil, repo, redisDB, nil)
		err := service.UpdatePsw(ctx, uid, oldPsw, newPsw, code)
		assert.Error(t, err)
	})

	t.Run("new psw is invalid", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		repo.On("UpdatePsw", ctx, uid, mock.Anything).Return(nil)
		redisDB := new(mockRedis)
		redisDB.On("Get", ctx, mock.Anything).Return(code, nil)
		service := NewService(nil, repo, redisDB, nil)
		err := service.UpdatePsw(ctx, uid, oldPsw, newPsw, code)
		assert.NoError(t, err)
	})
}

func TestService_SendSecurityCodeForUpdateEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("not found user", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(nil, fmt.Errorf("error"))

		service := NewService(nil, repo, nil, nil)
		err := service.SendSecurityCodeForUpdateEmail(ctx, uid)
		assert.Error(t, err)
	})

	t.Run("set code in redis", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("SetNX", ctx, mock.Anything, mock.Anything, 5*time.Minute).Return(false, fmt.Errorf("error"))
		service := NewService(nil, repo, redisDB, nil)
		err := service.SendSecurityCodeForUpdateEmail(ctx, uid)
		assert.Error(t, err)
	})

	t.Run("code already exist in redis", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("SetNX", ctx, mock.Anything, mock.Anything, 5*time.Minute).Return(false, nil)
		service := NewService(nil, repo, redisDB, nil)
		err := service.SendSecurityCodeForUpdateEmail(ctx, uid)
		assert.Error(t, err)
	})

	t.Run("sending code is failed", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("SetNX", ctx, mock.Anything, mock.Anything, 5*time.Minute).Return(true, nil)
		emailMockSvc := new(mockEmailSvc)
		emailMockSvc.On("SendEmailForUpdateEmail", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("error"))
		service := NewService(nil, repo, redisDB, emailMockSvc)
		err := service.SendSecurityCodeForUpdateEmail(ctx, uid)
		assert.Error(t, err)
	})

	t.Run("sending code is done", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("SetNX", ctx, mock.Anything, mock.Anything, 5*time.Minute).Return(true, nil)
		emailMockSvc := new(mockEmailSvc)
		emailMockSvc.On("SendEmailForUpdateEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		service := NewService(nil, repo, redisDB, emailMockSvc)
		err := service.SendSecurityCodeForUpdateEmail(ctx, uid)
		assert.NoError(t, err)
	})
}

func TestService_UpdateEmail(t *testing.T) {
	ctx := context.Background()
	newEmail := "new_test@test.com"
	code := "1234"

	t.Run("new email is the same as old", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		service := NewService(nil, repo, nil, nil)
		wrongNewEmail := "test@test.com"
		err := service.UpdateEmail(ctx, uid, wrongNewEmail, code)
		assert.Error(t, err)
	})

	t.Run("new email is invalid", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		service := NewService(nil, repo, nil, nil)
		wrongNewEmail := "new_test@test"
		err := service.UpdateEmail(ctx, uid, wrongNewEmail, code)
		assert.Error(t, err)
	})

	t.Run("not found redis code", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("Get", ctx, mock.Anything).Return("", fmt.Errorf("error"))
		service := NewService(nil, repo, redisDB, nil)
		err := service.UpdateEmail(ctx, uid, newEmail, code)
		assert.Error(t, err)
	})

	t.Run("redis code is different", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		redisDB := new(mockRedis)
		redisDB.On("Get", ctx, mock.Anything).Return("1233", nil)
		service := NewService(nil, repo, redisDB, nil)
		err := service.UpdateEmail(ctx, uid, newEmail, code)
		assert.Error(t, err)
	})

	t.Run("update email is failed", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		repo.On("UpdateEmail", ctx, uid, mock.Anything).Return(fmt.Errorf("error"))
		redisDB := new(mockRedis)
		redisDB.On("Get", ctx, mock.Anything).Return("1234", nil)
		service := NewService(nil, repo, redisDB, nil)
		err := service.UpdateEmail(ctx, uid, newEmail, code)
		assert.Error(t, err)
	})

	t.Run("update email is done", func(t *testing.T) {
		repo := new(mockUserRepo)
		repo.On("GetUserByID", ctx, uid).Return(entUser, nil)
		repo.On("UpdateEmail", ctx, uid, mock.Anything).Return(nil)
		redisDB := new(mockRedis)
		redisDB.On("Get", ctx, mock.Anything).Return("1234", nil)
		service := NewService(nil, repo, redisDB, nil)
		err := service.UpdateEmail(ctx, uid, newEmail, code)
		assert.NoError(t, err)
	})
}
