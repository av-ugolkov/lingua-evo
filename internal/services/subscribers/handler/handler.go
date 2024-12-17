package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/services/subscribers"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	paramsSubscriberID = "subscriber_id"
)

type (
	SubscribeRs struct {
		ID uuid.UUID `json:"id"`
	}
)

type Handler struct {
	subscribersSvc *subscribers.Service
}

func Create(r *fiber.App, userSvc *subscribers.Service) {
	h := newHandler(userSvc)

	r.Post(handler.Subscribe, middleware.Auth(h.subscribe))
	r.Post(handler.Unsubscribe, middleware.Auth(h.unsubscribe))
	r.Get(handler.CheckSubscriber, middleware.Auth(h.checkSubscriber))
}

func newHandler(subscribersSvc *subscribers.Service) *Handler {
	return &Handler{
		subscribersSvc: subscribersSvc,
	}
}

func (h *Handler) subscribe(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized,
			fmt.Sprintf("subscribers.delivery.Handler.subscribe: %v", err))
	}

	var data SubscribeRs
	err = c.BodyParser(&data)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgInternal)
	}

	err = h.subscribersSvc.Subscribe(ctx, uid, data.ID)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgInternal)
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) unsubscribe(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, msgerr.ErrMsgUnauthorized)
	}

	var data SubscribeRs
	err = c.BodyParser(&data)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgInternal)
	}

	err = h.subscribersSvc.Unsubscribe(ctx, uid, data.ID)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, msgerr.ErrMsgInternal)
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) checkSubscriber(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized,
			fmt.Sprintf("subscribers.delivery.Handler.checkSubscriber: %v", err))
	}

	subID, err := uuid.Parse(c.Query(paramsSubscriberID))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("subscribers.delivery.Handler.checkSubscriber: %v", err))
	}

	isSubscriber, err := h.subscribersSvc.Check(ctx, uid, subID)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("subscribers.delivery.Handler.checkSubscriber: %v", err))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"is_subscriber": isSubscriber})
}
