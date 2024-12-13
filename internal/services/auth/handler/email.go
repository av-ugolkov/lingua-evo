package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	ginext "github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

func (h *Handler) signIn(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	authorization, err := c.GetHeaderAuthorization(ginext.AuthTypeBasic)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("auth.handler.Handler.signIn: %w", err)
	}

	var data CreateSessionRq
	err = decodeBasicAuth(authorization, &data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("auth.handler.Handler.signIn: %v", err)
	}
	var fingerprint string
	if fingerprint = c.GetHeader(ginext.Fingerprint); fingerprint == runtime.EmptyString {
		return http.StatusBadRequest, nil,
			fmt.Errorf("auth.handler.Handler.signIn: fingerprint not found")
	}
	data.Fingerprint = fingerprint

	refreshTokenID := uuid.New()
	tokens, err := h.authSvc.SignIn(ctx, data.User, data.Password, data.Fingerprint, refreshTokenID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFoundUser) ||
			errors.Is(err, auth.ErrWrongPassword):
			return http.StatusNotFound, nil,
				msgerr.New(fmt.Errorf("auth.handler.Handler.signIn: %w", err),
					"User doesn't exist or password is wrong")
		default:
			return http.StatusInternalServerError, nil,
				msgerr.New(fmt.Errorf("auth.handler.Handler.signIn: %w", err),
					msgerr.ErrMsgInternal)
		}
	}

	sessionRs := &CreateSessionRs{
		AccessToken: tokens.AccessToken,
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second

	c.SetCookieRefreshToken(tokens.RefreshToken, duration)
	return http.StatusOK, sessionRs, nil
}

func (h *Handler) signUp(c *ginext.Context) (int, any, error) {
	var data CreateUserRq
	err := c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(
				fmt.Errorf("auth.handler.Handler.signUp: %v", err),
				msgerr.ErrMsgBadRequest)

	}

	if !utils.IsPasswordValid(data.Password) {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("auth.handler.Handler.signUp: invalid password"),
				"Invalid password")
	}

	if !utils.IsEmailValid(data.Email) {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("auth.handler.Handler.signUp: invalid email"),
				msgerr.ErrMsgBadEmail)
	}

	fingerprint := c.GetHeader(ginext.Fingerprint)
	if fingerprint == runtime.EmptyString {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.handler.Handler.signUp: fingerprimt is empty")
	}

	uid, err := h.authSvc.SignUp(c.Request.Context(), entity.User{
		Nickname: runtime.GenerateNickname(),
		Password: data.Password,
		Email:    data.Email,
		Role:     runtime.User,
		Code:     data.Code,
	}, fingerprint)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("auth.handler.Handler.signUp: %w", err)
	}

	createUserRs := &CreateUserRs{
		UserID: uid,
	}

	return http.StatusCreated, createUserRs, nil
}

func (h *Handler) sendCode(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var data CreateCodeRq
	err := c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("auth.handler.Handler.sendCode: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	if !utils.IsEmailValid(data.Email) {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("auth.handler.Handler.sendCode - email format is invalid"),
				msgerr.ErrMsgBadEmail)
	}

	fingerprint := c.GetHeader(ginext.Fingerprint)
	if fingerprint == runtime.EmptyString {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.handler.Handler.sendCode: fingerprimt is empty")
	}

	err = h.authSvc.CreateCode(ctx, data.Email, fingerprint)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("auth.handler.Handler.sendCode: %v", err),
				msgerr.ErrMsgInternal)
	}

	return http.StatusOK, nil, nil
}

func decodeBasicAuth(basicToken string, data *CreateSessionRq) error {
	base, err := base64.StdEncoding.DecodeString(basicToken)
	if err != nil {
		return fmt.Errorf("auth.handler.decodeBasicAuth: %v", err)
	}
	authData := strings.Split(string(base), ":")
	if len(authData) != 2 {
		return fmt.Errorf("auth.handler.decodeBasicAuth: invalid auth data")
	}

	data.User = authData[0]
	data.Password = authData[1]

	return nil
}
