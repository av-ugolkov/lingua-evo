package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/services/access"
)

type (
	Access struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	}

	AccessesRs struct {
		Accesses []Access `json:"accesses"`
	}
)

type Handler struct {
	accessSvc *access.Service
}

func Create(r *ginext.Engine, accessSvc *access.Service) {
	h := newHandler(accessSvc)

	r.GET(handler.Accesses, h.getAccesses)
}

func newHandler(accessSvc *access.Service) *Handler {
	return &Handler{
		accessSvc: accessSvc,
	}
}

func (h *Handler) getAccesses(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	accesses, err := h.accessSvc.GetAccesses(ctx)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("access.delivery.Handler.getAccesses: %v", err)
	}

	accessesRs := make([]Access, 0, len(accesses))
	for _, access := range accesses {
		accessesRs = append(accessesRs, Access{
			ID:   access.ID,
			Type: access.Type,
			Name: access.Name,
		})
	}

	return http.StatusOK, accessesRs, nil
}
