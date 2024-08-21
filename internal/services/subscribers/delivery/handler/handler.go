package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
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
	r.POST(handler.Subscribe, middleware.Auth(h.subscribe))
	r.POST(handler.Unsubscribe, middleware.Auth(h.unsubscribe))
	r.GET(handler.CheckSubscriber, middleware.Auth(h.checkSubscriber))
}

func (h *Handler) subscribe(c *gin.Context) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("subscribers.delivery.handler.Handler.subscribe - get user id: %v", err))
		return
	}

	var data SubscribeRs
	err = c.Bind(&data)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("subscribers.delivery.handler.Handler.subscribe - check body: %v", err))
		return
	}

	err = h.subscribersSvc.Subscribe(ctx, uid, data.ID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("subscribers.delivery.handler.Handler.subscribe - subscribe: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) unsubscribe(c *gin.Context) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("subscribers.delivery.handler.Handler.unsubscribe - get user id: %v", err))
		return
	}

	var data SubscribeRs
	err = c.Bind(&data)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("subscribers.delivery.handler.Handler.unsubscribe - check body: %v", err))
		return
	}

	err = h.subscribersSvc.Unsubscribe(ctx, uid, data.ID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("subscribers.delivery.handler.Handler.unsubscribe - subscribe: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) checkSubscriber(c *gin.Context) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("subscribers.delivery.handler.Handler.checkSubscriber - get user id: %v", err))
		return
	}

	subID, err := ginExt.GetQueryUUID(c, paramsSubscriberID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("subscribers.delivery.handler.Handler.checkSubscriber - get query [sub_id]: %v", err))
		return
	}

	isSubscriber, err := h.subscribersSvc.Check(ctx, uid, subID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("subscribers.delivery.handler.Handler.checkSubscriber - check: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_subscriber": isSubscriber,
	})
}
