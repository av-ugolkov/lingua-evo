package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/router"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth/dto"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initGoogleHandler(r *fiber.App) {
	r.Get(handler.GoogleAuth, h.googleAuthUrl)
	r.Post(handler.GoogleAuth, h.googleAuth)
}

func (h *Handler) googleAuthUrl(c *fiber.Ctx) error {
	url := h.authSvc.GoogleAuthUrl()
	return c.Status(http.StatusOK).JSON(fiber.Map{"url": url})
}

func (h *Handler) googleAuth(c *fiber.Ctx) error {
	ctx := c.Context()

	var data dto.GoogleAuthCode
	err := c.BodyParser(&data)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgBadRequest)
	}

	fingerprint := c.GetReqHeaders()[router.Fingerprint]
	if fingerprint[0] == runtime.EmptyString {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("auth.handler.Handler.googleAuth: fingerprimt is empty"))
	}

	tokens, err := h.authSvc.AuthByGoogle(ctx, data.Code, fingerprint[0])
	if err != nil {
		var msgErr *msgerr.Error
		switch {
		case errors.Is(err, entity.ErrNotFoundUser) ||
			errors.Is(err, auth.ErrWrongPassword):
			return fiber.NewError(http.StatusBadRequest,
				"User doesn't exist or password is wrong")
		case errors.As(err, &msgErr):
			return fiber.NewError(http.StatusBadRequest, msgErr.Msg)
		default:
			return fiber.NewError(http.StatusInternalServerError,
				"Sorry! We can't sign you in. Please try again.")
		}
	}

	sessionRs := &dto.CreateSessionRs{
		AccessToken: tokens.AccessToken,
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second

	c.Cookie(&fiber.Cookie{
		Name:     router.RefreshToken,
		Value:    tokens.RefreshToken,
		MaxAge:   int(duration.Seconds()),
		Path:     router.CookiePathAuth,
		Secure:   true,
		HTTPOnly: true,
	})
	return c.Status(http.StatusOK).JSON(sessionRs)
}
