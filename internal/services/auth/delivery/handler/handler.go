package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/pkg/http/gin_extension"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
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

func Create(r *gin.Engine, authSvc *auth.Service) {
	h := newHandler(authSvc)
	h.register(r)
}

func newHandler(authSvc *auth.Service) *Handler {
	return &Handler{
		authSvc: authSvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.POST(delivery.SignIn, h.signIn)
	r.GET(delivery.Refresh, h.refresh)
	r.GET(delivery.SignOut, middleware.Auth(h.signOut))
	r.POST(delivery.SendCode, h.sendCode)
}

func (h *Handler) signIn(c *gin.Context) {
	ctx := c.Request.Context()
	authorization, err := ginExt.GetHeaderAuthorization(c, ginExt.AuthTypeBasic)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("auth.delivery.Handler.signin: %w", err))
		return
	}

	var data CreateSessionRq
	err = decodeBasicAuth(authorization, &data)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("auth.delivery.Handler.signin - check body: %v", err))
		return
	}
	var fingerprint string
	if fingerprint = c.GetHeader(ginExt.Fingerprint); fingerprint == runtime.EmptyString {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("auth.delivery.Handler.signin: fingerprint not found"))
		return
	}

	refreshTokenID := uuid.New()
	tokens, err := h.authSvc.SignIn(ctx, data.User, data.Password, data.Fingerprint, refreshTokenID)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("auth.delivery.Handler.signin - create session: %v", err))
		return
	}

	sessionRs := &CreateSessionRs{
		AccessToken: tokens.AccessToken,
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second

	ginExt.SetCookieRefreshToken(c, tokens.RefreshToken, duration)
	c.JSON(http.StatusOK, sessionRs)
}

func (h *Handler) refresh(c *gin.Context) {
	ctx := c.Request.Context()
	refreshToken, err := c.Cookie(ginExt.RefreshToken)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get cookie: %v", err))
		return
	}

	refreshID, err := uuid.Parse(refreshToken)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("auth.delivery.Handler.refresh - parse: %v", err))
		return
	}
	var fingerprint string
	if fingerprint = c.GetHeader(ginExt.Fingerprint); fingerprint == runtime.EmptyString {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("auth.delivery.Handler.refresh - get fingerprint: %v", err))
		return
	}

	tokenID := uuid.New()
	tokens, err := h.authSvc.RefreshSessionToken(ctx, tokenID, refreshID, fingerprint)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("auth.delivery.Handler.refresh - refresh token: %v", err))
		return
	}

	sessionRs := &CreateSessionRs{
		AccessToken: tokens.AccessToken,
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	ginExt.SetCookieRefreshToken(c, tokens.RefreshToken, duration)
	c.JSON(http.StatusOK, sessionRs)
}

func (h *Handler) signOut(c *gin.Context) {
	ctx := c.Request.Context()
	refreshToken, err := c.Cookie(ginExt.RefreshToken)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - get cookie: %v", err))
		return
	}

	refreshID, err := uuid.Parse(refreshToken)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - parse: %v", err))
		return
	}

	var fingerprint string
	if fingerprint = c.GetHeader(ginExt.Fingerprint); fingerprint == runtime.EmptyString {
		ginExt.SendError(c, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - get fingerprint: %v", err))
		return
	}

	err = h.authSvc.SignOut(ctx, refreshID, fingerprint)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - logout: %v", err))
		return
	}

	ginExt.DeleteCookie(c, ginExt.RefreshToken)
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) sendCode(c *gin.Context) {
	ctx := c.Request.Context()

	var data CreateCodeRq
	err := c.Bind(&data)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest, fmt.Errorf("auth.delivery.handler.Handler.sendCode - check body: %v", err))
		return
	}

	err = h.authSvc.CreateCode(ctx, data.Email)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.sendCode - create code: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func decodeBasicAuth(basicToken string, data *CreateSessionRq) error {
	base, err := base64.StdEncoding.DecodeString(basicToken)
	if err != nil {
		return fmt.Errorf("auth.delivery.decodeBasicAuth - decode base64: %v", err)
	}
	authData := strings.Split(string(base), ":")
	if len(authData) != 2 {
		return fmt.Errorf("auth.delivery.decodeBasicAuth - invalid auth data")
	}

	data.User = authData[0]
	data.Password = authData[1]

	return nil
}
