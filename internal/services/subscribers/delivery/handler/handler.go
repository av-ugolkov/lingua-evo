package handler

import (
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/services/subscribers"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	subscribersSvc *subscribers.Service
}

func Create(r *gin.Engine, userSvc *subscribers.Service) {
	h := newHandler(userSvc)
	h.register(r)
}

func newHandler(subscribersSvc *subscribers.Service) *Handler {
	return &Handler{
		subscribersSvc: subscribersSvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.GET(handler.Subscribers, middleware.Auth(h.getRespondents))
}

func (h *Handler) getRespondents(c *gin.Context) {

}
