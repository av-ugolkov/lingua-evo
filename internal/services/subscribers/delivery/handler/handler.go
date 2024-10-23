package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/services/subscribers"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
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

func Create(r *ginext.Engine, userSvc *subscribers.Service) {
	h := newHandler(userSvc)

	r.POST(handler.Subscribe, middleware.Auth(h.subscribe))
	r.POST(handler.Unsubscribe, middleware.Auth(h.unsubscribe))
	r.GET(handler.CheckSubscriber, middleware.Auth(h.checkSubscriber))
}

func newHandler(subscribersSvc *subscribers.Service) *Handler {
	return &Handler{
		subscribersSvc: subscribersSvc,
	}
}

func (h *Handler) subscribe(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("subscribers.delivery.handler.Handler.subscribe: %v", err)
	}

	var data SubscribeRs
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("subscribers.delivery.handler.Handler.subscribe: %v", err)
	}

	err = h.subscribersSvc.Subscribe(ctx, uid, data.ID)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("subscribers.delivery.handler.Handler.subscribe: %v", err)
	}

	return http.StatusOK, gin.H{}, nil
}

func (h *Handler) unsubscribe(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("subscribers.delivery.handler.Handler.unsubscribe: %v", err)
	}

	var data SubscribeRs
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("subscribers.delivery.handler.Handler.unsubscribe: %v", err)
	}

	err = h.subscribersSvc.Unsubscribe(ctx, uid, data.ID)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("subscribers.delivery.handler.Handler.unsubscribe: %v", err)
	}

	return http.StatusOK, gin.H{}, nil
}

func (h *Handler) checkSubscriber(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("subscribers.delivery.handler.Handler.checkSubscriber: %v", err)
	}

	subID, err := c.GetQueryUUID(paramsSubscriberID)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("subscribers.delivery.handler.Handler.checkSubscriber: %v", err)
	}

	isSubscriber, err := h.subscribersSvc.Check(ctx, uid, subID)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("subscribers.delivery.handler.Handler.checkSubscriber: %v", err)
	}

	return http.StatusOK, gin.H{
		"is_subscriber": isSubscriber,
	}, nil
}
