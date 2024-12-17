package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/services/notifications"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	paramsUserID  string = "user_id"
	paramsVocabID string = "vocab_id"
)

type Handler struct {
	notificationsSvc *notifications.Service
}

func Create(g *fiber.App, notificationsSvc *notifications.Service) {
	h := newHandler(notificationsSvc)

	g.Post(handler.NotificationVocab, h.setNotificationVocab)
}

func newHandler(notificationsSvc *notifications.Service) *Handler {
	return &Handler{
		notificationsSvc: notificationsSvc,
	}
}

func (h *Handler) setNotificationVocab(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Query(paramsUserID))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("notifications.delivery.Handler.setNotificationVocab: %v", err))
	}

	vid, err := uuid.Parse(c.Query(paramsVocabID))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("notifications.delivery.Handler.setNotificationVocab: %v", err))
	}

	ok, err := h.notificationsSvc.SetVocabNotification(c.Context(), uid, vid)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("notifications.delivery.Handler.setNotificationVocab: %w", err))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"notification": ok})
}
