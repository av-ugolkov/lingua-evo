package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/services/notifications"

	"github.com/gin-gonic/gin"
)

const (
	paramsUserID  string = "user_id"
	paramsVocabID string = "vocab_id"
)

type Handler struct {
	notificationsSvc *notifications.Service
}

func Create(g *ginext.Engine, notificationsSvc *notifications.Service) {
	h := newHandler(notificationsSvc)

	g.POST(handler.NotificationVocab, h.setNotificationVocab)
}

func newHandler(notificationsSvc *notifications.Service) *Handler {
	return &Handler{
		notificationsSvc: notificationsSvc,
	}
}

func (h *Handler) setNotificationVocab(c *ginext.Context) (int, any, error) {
	uid, err := c.GetQueryUUID(paramsUserID)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("notifications.delivery.Handler.setNotificationVocab: %v", err)
	}

	vid, err := c.GetQueryUUID(paramsVocabID)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("notifications.delivery.Handler.setNotificationVocab: %v", err)
	}

	ok, err := h.notificationsSvc.SetVocabNotification(c.Request.Context(), uid, vid)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("notifications.delivery.Handler.setNotificationVocab: %w", err)
	}

	return http.StatusOK, gin.H{"notification": ok}, nil
}
