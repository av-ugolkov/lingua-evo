package handler

import (
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
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

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	vid, err := uuid.Parse(c.Query(paramsID))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err))
	}

	accessVocab, err := h.vocabSvc.GetAccessForUser(ctx, uid, vid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.Status(http.StatusOK).JSON(fext.D(fiber.Map{"access": accessVocab}))
}

func (h *Handler) addAccessForUser(c *fiber.Ctx) error {
	ctx := c.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.BodyParser(&vocabAccessRq)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err))
	}

	err = h.vocabSvc.AddAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
		Status:  vocabAccessRq.Status,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) removeAccessForUser(c *fiber.Ctx) error {
	ctx := c.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.BodyParser(&vocabAccessRq)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err))
	}

	err = h.vocabSvc.RemoveAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) updateAccessForUser(c *fiber.Ctx) error {
	ctx := c.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.BodyParser(&vocabAccessRq)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err))
	}

	err = h.vocabSvc.UpdateAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
		Status:  vocabAccessRq.Status,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	return c.SendStatus(http.StatusOK)
}
