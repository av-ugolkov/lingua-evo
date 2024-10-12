package handler

import (
	"fmt"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/services/notifications"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	paramsUserID  string = "user_id"
	paramsVocabID string = "vocab_id"
)

type Handler struct {
	notificationsSvc *notifications.Service
}

func Create(g *gin.Engine, notificationsSvc *notifications.Service) {
	h := newHandler(notificationsSvc)
	h.register(g)
}

func newHandler(notificationsSvc *notifications.Service) *Handler {
	return &Handler{
		notificationsSvc: notificationsSvc,
	}
}

func (h *Handler) register(g *gin.Engine) {
	g.GET(handler.NotificationVocab, h.getNotificationVocab)
	g.POST(handler.NotificationVocab, h.setNotificationVocab)
}

func (h *Handler) getNotificationVocab(c *gin.Context) {
	uid, err := ginExt.GetQueryUUID(c, paramsUserID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("notifications.delivery.Handler.getNotificationVocab - get query [user_id]: %w", err))
		return
	}

	vid, err := ginExt.GetQueryUUID(c, paramsVocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("notifications.delivery.Handler.getNotificationVocab - get query [vocab_id]: %w", err))
		return
	}

	ok, err := h.notificationsSvc.GetVocabNotification(c.Request.Context(), uid, vid)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("notifications.delivery.Handler.getNotificationVocab - get notification: %w", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification": ok})
}

func (h *Handler) setNotificationVocab(c *gin.Context) {
	uid, err := ginExt.GetQueryUUID(c, paramsUserID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("notifications.delivery.Handler.setNotificationVocab - get query [user_id]: %w", err))
		return
	}

	vid, err := ginExt.GetQueryUUID(c, paramsVocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("notifications.delivery.Handler.setNotificationVocab - get query [vocab_id]: %w", err))
		return
	}

	ok, err := h.notificationsSvc.SetVocabNotification(c.Request.Context(), uid, vid)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("notifications.delivery.Handler.setNotificationVocab - set notification: %w", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification": ok})
}
