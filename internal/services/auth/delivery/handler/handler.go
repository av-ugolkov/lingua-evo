package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	CreateSessionRq struct {
		User        string `json:"user"`
		Password    string `json:"password"`
		Fingerprint string `json:"fingerprint"`
	}

	CreateCodeRq struct {
		Email string `json:"email"`
	}

	CreateSessionRs struct {
		AccessToken string `json:"access_token"`
	}
)

type Handler struct {
	authSvc *auth.Service
}

func Create(r *ginext.Engine, authSvc *auth.Service) {
	h := newHandler(authSvc)

	r.POST(handler.SignIn, h.signIn)
	r.GET(handler.Refresh, h.refresh)
	r.POST(handler.SignOut, middleware.Auth(h.signOut))
	r.POST(handler.SendCode, h.sendCode)
}

func newHandler(authSvc *auth.Service) *Handler {
	return &Handler{
		authSvc: authSvc,
	}
}

func (h *Handler) signIn(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	authorization, err := c.GetHeaderAuthorization(ginext.AuthTypeBasic)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("auth.delivery.Handler.signin: %w", err)
	}

	var data CreateSessionRq
	err = decodeBasicAuth(authorization, &data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("auth.delivery.Handler.signin: %v", err)
	}
	var fingerprint string
	if fingerprint = c.GetHeader(ginext.Fingerprint); fingerprint == runtime.EmptyString {
		return http.StatusBadRequest, nil,
			fmt.Errorf("auth.delivery.Handler.signin: fingerprint not found")
	}
	data.Fingerprint = fingerprint

	refreshTokenID := uuid.New()
	tokens, err := h.authSvc.SignIn(ctx, data.User, data.Password, data.Fingerprint, refreshTokenID)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("auth.delivery.Handler.signin: %v", err)
	}

	sessionRs := &CreateSessionRs{
		AccessToken: tokens.AccessToken,
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second

	c.SetCookieRefreshToken(tokens.RefreshToken, duration)
	return http.StatusOK, sessionRs, nil
}

func (h *Handler) refresh(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	refreshToken, err := c.Cookie(ginext.RefreshToken)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("auth.delivery.Handler.refresh: %v", err)
	}
	if refreshToken == runtime.EmptyString {
		return http.StatusBadRequest, nil, fmt.Errorf("auth.delivery.Handler.refresh - refresh token not found")
	}

	refreshID, err := uuid.Parse(refreshToken)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("auth.delivery.Handler.refresh: %v", err)
	}
	var fingerprint string
	if fingerprint = c.GetHeader(ginext.Fingerprint); fingerprint == runtime.EmptyString {
		return http.StatusBadRequest, nil,
			fmt.Errorf("auth.delivery.Handler.refresh: fingerprint is empty")
	}

	tokenID := uuid.New()
	tokens, err := h.authSvc.RefreshSessionToken(ctx, tokenID, refreshID, fingerprint)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("auth.delivery.Handler.refresh: %v", err)
	}

	sessionRs := &CreateSessionRs{
		AccessToken: tokens.AccessToken,
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	c.SetCookieRefreshToken(tokens.RefreshToken, duration)
	return http.StatusOK, sessionRs, nil
}

func (h *Handler) signOut(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	refreshToken, err := c.Cookie(ginext.RefreshToken)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.delivery.Handler.signOut: %v", err)
	}

	refreshID, err := uuid.Parse(refreshToken)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.delivery.Handler.signOut: %v", err)
	}

	var fingerprint string
	if fingerprint = c.GetHeader(ginext.Fingerprint); fingerprint == runtime.EmptyString {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.delivery.Handler.signOut: fingerprimt is empty")
	}

	err = h.authSvc.SignOut(ctx, refreshID, fingerprint)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.delivery.Handler.signOut: %v", err)
	}

	c.DeleteCookie(ginext.RefreshToken)
	return http.StatusOK, gin.H{}, nil
}

func (h *Handler) sendCode(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var data CreateCodeRq
	err := c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("auth.delivery.Handler.sendCode: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	if !utils.IsEmailValid(data.Email) {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("auth.delivery.Handler.sendCode - email format is invalid"),
				msgerr.ErrMsgBadEmail)
	}

	err = h.authSvc.CreateCode(ctx, data.Email)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("auth.delivery.Handler.sendCode: %v", err),
				msgerr.ErrMsgInternal)
	}

	return http.StatusOK, gin.H{}, nil
}

func decodeBasicAuth(basicToken string, data *CreateSessionRq) error {
	base, err := base64.StdEncoding.DecodeString(basicToken)
	if err != nil {
		return fmt.Errorf("auth.delivery.decodeBasicAuth: %v", err)
	}
	authData := strings.Split(string(base), ":")
	if len(authData) != 2 {
		return fmt.Errorf("auth.delivery.decodeBasicAuth: invalid auth data")
	}

	data.User = authData[0]
	data.Password = authData[1]

	return nil
}
