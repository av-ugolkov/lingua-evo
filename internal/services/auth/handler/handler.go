package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/router"
	auth "github.com/av-ugolkov/lingua-evo/internal/services/auth/service"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	authSvc *auth.Service
}

func Create(r *fiber.App, authSvc *auth.Service) {
	h := &Handler{
		authSvc: authSvc,
	}

	r.Get(handler.Refresh, h.refresh)
	r.Post(handler.SignOut, middleware.Auth(h.signOut))

	h.initEmailHandler(r)
	h.initGoogleHandler(r)
}

func (h *Handler) refresh(c *fiber.Ctx) error {
	ctx := c.Context()
	var err error
	defer func() {
		if err != nil {
			c.ClearCookie(router.RefreshToken, router.CookiePathAuth)
		}
	}()

	refreshToken := c.Cookies(router.RefreshToken)
	if refreshToken == runtime.EmptyString {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("auth.handler.Handler.refresh - refresh token not found"))
	}

	fingerprint := c.GetReqHeaders()[router.Fingerprint]
	if fingerprint[0] == runtime.EmptyString {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("auth.handler.Handler.refresh: fingerprint is empty"))
	}

	uid, err := uuid.Parse(c.Query("uid"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("auth.handler.Handler.refresh: %v", err))
	}

	sessionRs, err := h.authSvc.RefreshSessionToken(ctx, uid, refreshToken, fingerprint[0])
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("auth.handler.Handler.refresh: %v", err))
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	c.Cookie(&fiber.Cookie{
		Name:     router.RefreshToken,
		Value:    refreshToken,
		MaxAge:   int(duration.Seconds()),
		Path:     router.CookiePathAuth,
		Secure:   true,
		HTTPOnly: true,
	})

	return c.Status(http.StatusOK).JSON(sessionRs)
}

func (h *Handler) signOut(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		c.ClearCookie(router.RefreshToken, router.CookiePathAuth)
		return fiber.NewError(http.StatusInternalServerError,
			msgerr.ErrMsgUnauthorized)
	}

	refreshToken := c.Cookies(router.RefreshToken)
	if refreshToken == runtime.EmptyString {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("auth.handler.Handler.signOut: %v", err))
	}

	fingerprint := c.GetReqHeaders()[router.Fingerprint]
	if fingerprint[0] == runtime.EmptyString {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("auth.handler.Handler.signOut: fingerprimt is empty"))
	}

	err = h.authSvc.SignOut(ctx, uid, refreshToken, fingerprint[0])
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("auth.handler.Handler.signOut: %v", err))
	}

	c.ClearCookie(router.RefreshToken, router.CookiePathAuth)
	return c.SendStatus(http.StatusOK)
}
