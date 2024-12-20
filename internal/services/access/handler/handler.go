package handler

import (
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/services/access/dto"
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

	accesses, err := h.accessSvc.GetAccesses(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	accessesRs := dto.AccessesToDto(accesses)

	return c.Status(http.StatusOK).JSON(fext.D(accessesRs))
}
