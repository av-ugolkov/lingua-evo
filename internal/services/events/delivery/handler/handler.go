package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	ginext "github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/services/events/service"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

type (
	CountEventsRs struct {
		Count int `json:"count"`
	}

	userData struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Role        string    `json:"role"`
		LastVisitAt time.Time `json:"last_visit_at"`
	}

	EventsRs struct {
		ID        uuid.UUID `json:"id"`
		User      userData  `json:"user"`
		Msg       string    `json:"msg"`
		CreatedAt time.Time `json:"created_at"`
		Watched   bool      `json:"watched"`
	}
)

type Handler struct {
	eventsSvc *service.Service
}

func Create(r *ginext.Engine, eventsSvc *service.Service) {
	h := &Handler{
		eventsSvc: eventsSvc,
	}

	r.GET(handler.CountEvents, middleware.Auth(h.getCountEvents))
	r.GET(handler.Events, middleware.Auth(h.getEvents))
	r.POST(handler.MarkWatched, middleware.Auth(h.markEventAsWatched))
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

	eventsRs := make([]EventsRs, 0, len(events))
	for _, event := range events {
		eventsRs = append(eventsRs, EventsRs{
			ID: event.ID,
			User: userData{
				ID:          event.User.ID,
				Name:        event.User.Name,
				Role:        event.User.Role,
				LastVisitAt: event.User.LastVisitAt,
			},
			Msg:       event.Payload.String(),
			CreatedAt: event.CreatedAt,
			Watched:   event.Watched,
		})
	}

	return http.StatusOK, eventsRs, nil
}

func (h *Handler) markEventAsWatched(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("events.delivery.Handler.markEventAsWatched: %w", err),
				msgerr.ErrMsgUnauthorized)
	}

	eid, err := c.GetQueryUUID("event_id")
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("events.delivery.Handler.markEventAsWatched: %w", err),
				msgerr.ErrMsgInternal)
	}

	err = h.eventsSvc.ReadEvent(ctx, uid, eid)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("events.delivery.Handler.markEventAsWatched: %w", err),
				msgerr.ErrMsgInternal)
	}

	return http.StatusOK, gin.H{}, nil
}
