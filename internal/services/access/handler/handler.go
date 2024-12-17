package handler

import (
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/services/access/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	accessSvc *service.Service
}

func Create(r *fiber.App, accessSvc *service.Service) {
	h := &Handler{
		accessSvc: accessSvc,
	}

	r.Get(handler.Accesses, h.getAccesses)
}

func (h *Handler) getAccesses(c *fiber.Ctx) error {
	ctx := c.Context()

	accessesRs, err := h.accessSvc.GetAccessesDTO(ctx)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, msgerr.ErrMsgInternal)
	}

	return c.Status(http.StatusOK).JSON(accessesRs)
}
