package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
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
	return c.Status(http.StatusOK).JSON(fext.D(fiber.Map{"url": url}))
}

func (h *Handler) googleAuth(c *fiber.Ctx) error {
	ctx := c.Context()

	var data dto.GoogleAuthCode
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	fingerprint := c.GetReqHeaders()[router.HeaderFingerprint]
	if fingerprint[0] == runtime.EmptyString {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(
			fmt.Errorf("auth.handler.Handler.googleAuth: fingerprimt is empty")))
	}

	tokens, err := h.authSvc.AuthByGoogle(ctx, data.Code, fingerprint[0])
	if err != nil {
		var msgErr *msgerr.MsgErr
		switch {
		case errors.Is(err, entity.ErrNotFoundUser) ||
			errors.Is(err, auth.ErrWrongPassword):
			return c.Status(http.StatusBadRequest).JSON(fext.E(err,
				"User doesn't exist or password is wrong"))
		case errors.As(err, &msgErr):
			return c.Status(http.StatusBadRequest).JSON(fext.E(msgErr))
		default:
			return c.Status(http.StatusInternalServerError).JSON(fext.E(err,
				"Sorry! We can't sign you in. Please try again."))
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
	return c.Status(http.StatusOK).JSON(fext.D(sessionRs))
}
