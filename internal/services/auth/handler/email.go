package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/router"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth/dto"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) initEmailHandler(r *fiber.App) {
	r.Post(handler.SignIn, h.signIn)
	r.Post(handler.SignUp, h.signUp)
	r.Post(handler.SendCode, h.sendCode)
}

func (h *Handler) signIn(c *fiber.Ctx) error {
	ctx := c.Context()
	authorization, err := middleware.GetTokenAuth(c, router.AuthTypeBasic)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("auth.Handler.signIn: %м", err))
	}

	var data dto.CreateSessionRq
	err = decodeBasicAuth(authorization, &data)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("auth.Handler.signIn: %v", err))
	}
	fingerprint := c.GetReqHeaders()[router.Fingerprint]
	if fingerprint[0] == runtime.EmptyString {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("auth.Handler.signIn: fingerprint not found"))
	}
	data.Fingerprint = fingerprint[0]

	refreshTokenID := uuid.New()
	sessionRs, refreshToken, err := h.authSvc.SignIn(ctx, data.User, data.Password, data.Fingerprint, refreshTokenID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFoundUser) ||
			errors.Is(err, auth.ErrWrongPassword):
			return fiber.NewError(http.StatusNotFound, "User doesn't exist or password is wrong")
		default:
			return fiber.NewError(http.StatusInternalServerError, msgerr.ErrMsgInternal)
		}
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

func (h *Handler) signUp(c *fiber.Ctx) error {
	var data dto.CreateUserRq
	err := c.BodyParser(&data)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgBadRequest)
	}

	if !utils.IsPasswordValid(data.Password) {
		return fiber.NewError(http.StatusBadRequest, "Invalid password")
	}

	if !utils.IsEmailValid(data.Email) {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgBadEmail)
	}

	fingerprint := c.GetReqHeaders()[router.Fingerprint]
	if fingerprint[0] == runtime.EmptyString {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("auth.Handler.signUp: fingerprimt is empty"))
	}

	uid, err := h.authSvc.SignUp(c.Context(), entity.User{
		Nickname: runtime.GenerateNickname(),
		Password: data.Password,
		Email:    data.Email,
		Role:     runtime.User,
		Code:     data.Code,
	}, fingerprint[0])
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("auth.Handler.signUp: %v", err))
	}

	createUserRs := &dto.CreateUserRs{
		UserID: uid,
	}

	return c.Status(http.StatusCreated).JSON(createUserRs)
}

func (h *Handler) sendCode(c *fiber.Ctx) error {
	ctx := c.Context()

	var data dto.CreateCodeRq
	err := c.BodyParser(&data)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgBadRequest)
	}

	if !utils.IsEmailValid(data.Email) {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgBadEmail)
	}

	fingerprint := c.GetReqHeaders()[router.Fingerprint]
	if fingerprint[0] == runtime.EmptyString {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("auth.Handler.sendCode: fingerprimt is empty"))
	}

	err = h.authSvc.CreateCode(ctx, data.Email, fingerprint[0])
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, msgerr.ErrMsgInternal)
	}

	return c.SendStatus(http.StatusOK)
}

func decodeBasicAuth(basicToken string, data *dto.CreateSessionRq) error {
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
