package handler

import (
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/services/access"

	"github.com/gin-gonic/gin"
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

func Create(r *gin.Engine, accessSvc *access.Service) {
	h := newHandler(accessSvc)
	h.register(r)
}

func newHandler(accessSvc *access.Service) *Handler {
	return &Handler{
		accessSvc: accessSvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.GET(handler.Accesses, h.getAccesses)
}

func (h *Handler) getAccesses(c *gin.Context) {
	ctx := c.Request.Context()

	accesses, err := h.accessSvc.GetAccesses(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError, err)
		return
	}

	accessesRs := make([]Access, 0, len(accesses))
	for _, access := range accesses {
		accessesRs = append(accessesRs, Access{
			ID:   access.ID,
			Type: access.Type,
			Name: access.Name,
		})
	}

	c.JSON(http.StatusOK, accessesRs)
}
