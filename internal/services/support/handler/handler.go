package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
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
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgBadRequest)
	}

	if !utils.IsEmailValid(data.Email) {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgBadEmail)
	}

	if len(data.Name) > 100 {
		data.Name = fmt.Sprintf("%s...", data.Name[:100])
	}

	if len(data.Message) == 0 {
		return fiber.NewError(http.StatusBadRequest, "Message is empty")
	}

	if len(data.Message) > 500 {
		return fiber.NewError(http.StatusBadRequest, "Message is too long")
	}

	err = h.supportSvc.SendRequest(ctx, support.SupportRequest{
		Email:   data.Email,
		Name:    data.Name,
		Type:    data.Type,
		Message: data.Message,
	})
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Something went wrong. Try a bit later!")
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"msg": "Thank you for your message! We will try to respond as soon as possible."})
}
