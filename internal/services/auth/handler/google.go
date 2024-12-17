package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginext "github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth/dto"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initGoogleHandler(r *ginext.Engine) {
	r.GET(handler.GoogleAuth, h.googleAuthUrl)
	r.POST(handler.GoogleAuth, h.googleAuth)
}

func (h *Handler) googleAuthUrl(c *ginext.Context) (int, any, error) {
	url := h.authSvc.GoogleAuthUrl()
	return http.StatusOK, gin.H{"url": url}, nil
}

func (h *Handler) googleAuth(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var data dto.GoogleAuthCode
	err := c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("auth.handler.Handler.googleAuth: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	fingerprint := c.GetHeader(ginext.Fingerprint)
	if fingerprint == runtime.EmptyString {
		return http.StatusInternalServerError, nil, fmt.Errorf("auth.handler.Handler.googleAuth: fingerprimt is empty")
	}

	tokens, err := h.authSvc.AuthByGoogle(ctx, data.Code, fingerprint)
	var e *msgerr.ApiError
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFoundUser) ||
			errors.Is(err, auth.ErrWrongPassword):
			return http.StatusBadRequest, nil,
				msgerr.New(fmt.Errorf("auth.handler.Handler.googleAuth: %w", err),
					"User doesn't exist or password is wrong")
		case errors.As(err, &e):
			return http.StatusBadRequest, nil, fmt.Errorf("auth.handler.Handler.googleAuth: %w", err)
		default:
			return http.StatusInternalServerError, nil,
				msgerr.New(fmt.Errorf("auth.handler.Handler.googleAuth: %w", err),
					msgerr.ErrMsgInternal)
		}
	}

	sessionRs := &dto.CreateSessionRs{
		AccessToken: tokens.AccessToken,
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second

	c.SetCookieRefreshToken(tokens.RefreshToken, duration)
	return http.StatusOK, sessionRs, nil
}
