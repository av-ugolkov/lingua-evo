package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	ginext "github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
)

func (h *Handler) googleAuthUrl(c *ginext.Context) (int, any, error) {
	url := h.authSvc.GoogleAuthUrl()
	return http.StatusOK, gin.H{"url": url}, nil
}

func (h *Handler) googleAuth(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var data GoogleAuthCode
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

	token, err := h.authSvc.AuthByGoogle(ctx, data.Code, fingerprint)
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

	sessionRs := &CreateSessionRs{
		AccessToken: token.AccessToken,
	}

	c.SetCookieRefreshToken(token.RefreshToken, time.Until(token.Expiry))
	return http.StatusOK, sessionRs, nil
}
