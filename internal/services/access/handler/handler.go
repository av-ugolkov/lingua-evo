package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/services/access/service"
)

type Handler struct {
	accessSvc *service.Service
}

func Create(r *ginext.Engine, accessSvc *service.Service) {
	h := &Handler{
		accessSvc: accessSvc,
	}

	r.GET(handler.Accesses, h.getAccesses)
}

func (h *Handler) getAccesses(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	accessesRs, err := h.accessSvc.GetAccessesDTO(ctx)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("access.delivery.Handler.getAccesses: %v", err)
	}

	return http.StatusOK, accessesRs, nil
}
