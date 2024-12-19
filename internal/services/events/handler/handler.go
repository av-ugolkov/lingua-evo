package handler

import (
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/services/events/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	CountEventsRs struct {
		Count int `json:"count"`
	}

	userData struct {
		ID        uuid.UUID `json:"id"`
		Nickname  string    `json:"nickname"`
		Role      string    `json:"role,omitempty"`
		VisitedAt time.Time `json:"visited_at,omitempty"`
	}

	EventsRs struct {
		ID        uuid.UUID      `json:"id"`
		User      userData       `json:"user"`
		Type      string         `json:"type"`
		Payload   map[string]any `json:"payload"`
		CreatedAt time.Time      `json:"created_at"`
		Watched   bool           `json:"watched"`
	}
)

type Handler struct {
	eventsSvc *service.Service
}

func Create(r *fiber.App, eventsSvc *service.Service) {
	h := &Handler{
		eventsSvc: eventsSvc,
	}

	r.Get(handler.CountEvents, middleware.Auth(h.getCountEvents))
	r.Get(handler.Events, middleware.Auth(h.getEvents))
	r.Post(handler.MarkWatched, middleware.Auth(h.markEventAsWatched))
}

func (h *Handler) getCountEvents(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	count, err := h.eventsSvc.GetCountEvents(ctx, uid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	return c.Status(http.StatusOK).JSON(fext.D(CountEventsRs{Count: count}))
}

func (h *Handler) getEvents(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	events, err := h.eventsSvc.GetEvents(ctx, uid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	eventsRs := make([]EventsRs, 0, len(events))
	for _, event := range events {
		eventsRs = append(eventsRs, EventsRs{
			ID: event.ID,
			User: userData{
				ID:        event.User.ID,
				Nickname:  event.User.Nickname,
				VisitedAt: event.User.VisitedAt,
			},
			Type:      string(event.Type),
			Payload:   event.PayloadToMap(),
			CreatedAt: event.CreatedAt,
			Watched:   event.Watched,
		})
	}

	return c.Status(http.StatusOK).JSON(fext.D(eventsRs))
}

func (h *Handler) markEventAsWatched(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	eid, err := uuid.Parse(c.Query("event_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	err = h.eventsSvc.ReadEvent(ctx, uid, eid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	return c.SendStatus(http.StatusOK)
}
