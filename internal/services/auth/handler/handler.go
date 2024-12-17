package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	auth "github.com/av-ugolkov/lingua-evo/internal/services/auth/service"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authSvc *auth.Service
}

func Create(r *ginext.Engine, authSvc *auth.Service) {
	h := &Handler{
		authSvc: authSvc,
	}

	r.GET(handler.Refresh, h.refresh)
	r.POST(handler.SignOut, middleware.Auth(h.signOut))

	h.initEmailHandler(r)
	h.initGoogleHandler(r)
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

	sessionRs, err := h.authSvc.RefreshSessionToken(ctx, uid, refreshToken, fingerprint)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("auth.handler.Handler.refresh: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	c.SetCookieRefreshToken(refreshToken, duration)

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
