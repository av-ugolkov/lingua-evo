package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	VocabularyAccessRq struct {
		ID      int           `json:"id,omitempty"`
		VocabID uuid.UUID     `json:"vocab_id,omitempty"`
		UserID  uuid.UUID     `json:"user_id,omitempty"`
		Status  access.Status `json:"access_edit,omitempty"`
	}
)

func (h *Handler) getAccessForUser(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized,
			fmt.Sprintf("vocabulary.Handler.getAccessForUser: %v", err))
	}

	vid, err := uuid.Parse(c.Query(paramsID))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.getAccessForUser: %v", err))
	}

	accessVocab, err := h.vocabSvc.GetAccessForUser(ctx, uid, vid)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.getAccessForUser: %v", err))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"access": accessVocab})
}

func (h *Handler) addAccessForUser(c *fiber.Ctx) error {
	ctx := c.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.BodyParser(&vocabAccessRq)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.addAccessForUser: %v", err))
	}

	err = h.vocabSvc.AddAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
		Status:  vocabAccessRq.Status,
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.addAccessForUser: %v", err))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) removeAccessForUser(c *fiber.Ctx) error {
	ctx := c.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.BodyParser(&vocabAccessRq)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.removeAccessForUser: %v", err))
	}

	err = h.vocabSvc.RemoveAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.removeAccessForUser: %v", err))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) updateAccessForUser(c *fiber.Ctx) error {
	ctx := c.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.BodyParser(&vocabAccessRq)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest,
			fmt.Sprintf("vocabulary.Handler.updateAccessForUser: %v", err))
	}

	err = h.vocabSvc.UpdateAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
		Status:  vocabAccessRq.Status,
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("vocabulary.Handler.updateAccessForUser: %v", err))
	}

	return c.SendStatus(http.StatusOK)
}
