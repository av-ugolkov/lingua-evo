package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	ginext "github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/services/events"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

type (
	CountEventsRs struct {
		Count int `json:"count"`
	}
	EventsRs struct {
		Events []events.Event `json:"events"`
	}
)

type Handler struct {
	eventsSvc *events.Service
}

func Create(r *ginext.Engine, eventsSvc *events.Service) {
	h := &Handler{
		eventsSvc: eventsSvc,
	}

	r.GET(handler.CountEvents, middleware.Auth(h.getCountEvents))
	r.GET(handler.Events, middleware.Auth(h.getEvents))
}

func (h *Handler) getCountEvents(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("events.delivery.Handler.getCountEvents: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	count, err := h.eventsSvc.GetCountEvents(ctx, uid)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("events.delivery.Handler.getCountEvents: %v", err),
				msgerr.ErrMsgInternal)
	}

	return http.StatusOK, CountEventsRs{Count: count}, nil
}

func (h *Handler) getEvents(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("events.delivery.Handler.getEvents: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	events, err := h.eventsSvc.GetEvents(ctx, uid)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("events.delivery.Handler.getEvents: %v", err),
				msgerr.ErrMsgInternal)
	}

	return http.StatusOK, EventsRs{Events: events}, nil
}
