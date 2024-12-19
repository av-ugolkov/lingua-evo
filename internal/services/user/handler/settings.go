package handler

import (
	"errors"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"

	"github.com/gofiber/fiber/v2"
)

type (
	UpdatePswRq struct {
		OldPsw string `json:"old_psw"`
		NewPsw string `json:"new_psw,omitempty"`
		Code   string `json:"code,omitempty"`
	}

	UpdateEmailRq struct {
		NewEmail string `json:"new_email"`
		Code     string `json:"code"`
	}

	UpdateNickname struct {
		Nickname string `json:"nickname"`
	}
)

func (h *Handler) initSettingsHandler(g *fiber.App) {
	g.Get(handler.AccountSettingsAccount, middleware.Auth(h.getSettingsAccount))
	g.Get(handler.AccountSettingsPersonalInfo, middleware.Auth(h.getSettingsPersonalInfo))
	g.Get(handler.AccountSettingsEmailNotif, middleware.Auth(h.getSettingsEmailNotif))
	g.Post(handler.AccountSettingsUpdatePswCode, middleware.Auth(h.updatePswSendCode))
	g.Post(handler.AccountSettingsUpdatePsw, middleware.Auth(h.updatePsw))
	g.Post(handler.AccountSettingsUpdateEmailCode, middleware.Auth(h.updateEmailSendCode))
	g.Post(handler.AccountSettingsUpdateEmail, middleware.Auth(h.updateEmail))
	g.Post(handler.AccountSettingsUpdateNickname, middleware.Auth(h.updateNickname))
}

func (h *Handler) getSettingsAccount(c *fiber.Ctx) error {
	ctx := c.Context()
	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	usr, err := h.userSvc.GetUserByID(ctx, uid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	userRs := UserRs{
		Nickname: usr.Nickname,
		Email:    usr.Email,
	}

	return c.Status(http.StatusOK).JSON(fext.D(userRs))
}

func (h *Handler) getSettingsPersonalInfo(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) getSettingsEmailNotif(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) updatePswSendCode(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data UpdatePswRq
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	ttl, err := h.userSvc.SendSecurityCodeForUpdatePsw(ctx, uid, data.OldPsw)
	switch {
	case errors.Is(err, entity.ErrDuplicateCode):
		return c.Status(http.StatusConflict).JSON(fext.DE(fiber.Map{"ttl": ttl}, err))
	case err != nil:
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) updatePsw(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data UpdatePswRq
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	err = h.userSvc.UpdatePsw(ctx, uid, data.OldPsw, data.NewPsw, data.Code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) updateEmailSendCode(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	ttl, err := h.userSvc.SendSecurityCodeForUpdateEmail(ctx, uid)
	switch {
	case errors.Is(err, entity.ErrDuplicateCode):
		return c.Status(http.StatusConflict).JSON(fext.DE(fiber.Map{"ttl": ttl}, err))
	case err != nil:
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.Status(http.StatusOK).JSON(fext.D(fiber.Map{
		"msg": "Could you check your email. We have sent you a code for updating your email",
	}))
}

func (h *Handler) updateEmail(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data UpdateEmailRq
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	err = h.userSvc.UpdateEmail(ctx, uid, data.NewEmail, data.Code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) updateNickname(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data UpdateNickname
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	err = h.userSvc.UpdateNickname(ctx, uid, data.Nickname)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.SendStatus(http.StatusOK)
}
