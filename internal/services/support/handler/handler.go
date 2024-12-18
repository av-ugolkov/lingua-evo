package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	"github.com/av-ugolkov/lingua-evo/internal/services/support"

	"github.com/gofiber/fiber/v2"
)

type (
	SupportRq struct {
		Email   string `json:"email"`
		Name    string `json:"name,omitempty"`
		Type    string `json:"type,omitempty"`
		Message string `json:"message"`
	}
)

type Handler struct {
	supportSvc *support.Service
}

func Create(r *fiber.App, supportSvc *support.Service) {
	h := &Handler{
		supportSvc: supportSvc,
	}

	r.Post(handler.SupportRequest, h.sendRequest)
}

func (h *Handler) sendRequest(c *fiber.Ctx) error {
	ctx := c.Context()

	var data SupportRq
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgBadRequest))
	}

	if !utils.IsEmailValid(data.Email) {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("support.Handler.sendRequest: email is invalid"),
			msgerr.ErrMsgBadEmail))
	}

	if len(data.Name) > 100 {
		data.Name = fmt.Sprintf("%s...", data.Name[:100])
	}

	if len(data.Message) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("support.Handler.sendRequest: message is empty")),
			"Message is empty")
	}

	if len(data.Message) > 500 {
		return c.Status(http.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("support.Handler.sendRequest: message is too long"),
			"Message is too long"))
	}

	err = h.supportSvc.SendRequest(ctx, support.SupportRequest{
		Email:   data.Email,
		Name:    data.Name,
		Type:    data.Type,
		Message: data.Message,
	})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, "Something went wrong. Try a bit later!"))
	}

	return c.Status(http.StatusCreated).JSON(fext.D(fiber.Map{
		"msg": "Thank you for your message! We will try to respond as soon as possible."}))
}
