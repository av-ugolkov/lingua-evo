package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/aes"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	auth "github.com/av-ugolkov/lingua-evo/internal/services/auth/service"
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

	CreateSessionRs struct {
		AccessToken string `json:"access_token"`
	}

	CreateUserRq struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Code     int    `json:"code"`
	}

	CreateUserRs struct {
		UserID uuid.UUID `json:"user_id"`
	}

	CreateCodeRq struct {
		Email string `json:"email"`
	}

	GoogleAuthCode struct {
		Code     string   `json:"code"`
		State    string   `json:"state"`
		Scope    []string `json:"scope"`
		Authuser int      `json:"authuser"`
		Prompt   string   `json:"prompt"`
	}
)

type Handler struct {
	authSvc *auth.Service
}

func Create(r *ginext.Engine, authSvc *auth.Service) {
	h := newHandler(authSvc)

	r.POST(handler.SignIn, h.signIn)
	r.POST(handler.SignUp, h.signUp)
	r.GET(handler.Refresh, h.refresh)
	r.POST(handler.SignOut, middleware.Auth(h.signOut))
	r.POST(handler.SendCode, h.sendCode)
	r.GET(handler.GoogleAuth, h.googleAuthUrl)
	r.POST(handler.GoogleAuth, h.googleAuth)
}

func newHandler(authSvc *auth.Service) *Handler {
	return &Handler{
		authSvc: authSvc,
	}
}

func (h *Handler) refresh(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	var err error
	defer func() {
		if err != nil {
			c.DeleteCookie(ginext.RefreshToken, "/auth")
		}
	}()

	refreshToken, err := c.Cookie(ginext.RefreshToken)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("auth.handler.Handler.refresh: %v", err)
	}
	if refreshToken == runtime.EmptyString {
		return http.StatusBadRequest, nil, fmt.Errorf("auth.handler.Handler.refresh - refresh token not found")
	}

	fingerprint := c.GetHeader(ginext.Fingerprint)
	if fingerprint == runtime.EmptyString {
		return http.StatusBadRequest, nil,
			fmt.Errorf("auth.handler.Handler.refresh: fingerprint is empty")
	}

	uid, err := c.GetQueryUUID("uid")
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("auth.handler.Handler.refresh: %v", err)
	}

	rt, err := aes.DecryptAES(refreshToken, config.GetConfig().AES.Key)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.handler.Handler.refresh: %v", err)
	}

	tokens, err := h.authSvc.RefreshSessionToken(ctx, uid, rt, fingerprint)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("auth.handler.Handler.refresh: %v", err)
	}

	sessionRs := &CreateSessionRs{
		AccessToken: tokens.AccessToken,
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second

	rt, err = aes.EncryptAES(tokens.RefreshToken, config.GetConfig().AES.Key)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.handler.Handler.refresh: %v", err)
	}

	c.SetCookieRefreshToken(rt, duration)
	return http.StatusOK, sessionRs, nil
}

func (h *Handler) signOut(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.handler.Handler.signOut: %v", err)
	}

	refreshToken, err := c.Cookie(ginext.RefreshToken)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.handler.Handler.signOut: %v", err)
	}

	fingerprint := c.GetHeader(ginext.Fingerprint)
	if fingerprint == runtime.EmptyString {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.handler.Handler.signOut: fingerprimt is empty")
	}

	err = h.authSvc.SignOut(ctx, uid, refreshToken, fingerprint)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.handler.Handler.signOut: %v", err)
	}

	c.DeleteCookie(ginext.RefreshToken, "/auth")
	return http.StatusOK, gin.H{}, nil
}
