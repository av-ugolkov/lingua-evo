package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	vocabAccess "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary_access"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	VocabularyAccessRq struct {
		ID         int       `json:"id,omitempty"`
		VocabID    uuid.UUID `json:"vocab_id,omitempty"`
		UserID     uuid.UUID `json:"user_id,omitempty"`
		AccessEdit bool      `json:"access_edit,omitempty"`
	}
)

type Handler struct {
	vocabAccessSvc *vocabAccess.Service
}

func Create(r *gin.Engine, vocabularySvc *vocabAccess.Service) {
	h := newHandler(vocabularySvc)
	h.register(r)
}

func newHandler(vocabAccessSvc *vocabAccess.Service) *Handler {
	return &Handler{
		vocabAccessSvc: vocabAccessSvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.POST(handler.VocabularyAccess, middleware.Auth(h.changeAccess))
	r.POST(handler.VocabularyAccessForUser, middleware.Auth(h.addAccessForUser))
	r.DELETE(handler.VocabularyAccessForUser, middleware.Auth(h.removeAccessForUser))
	r.PATCH(handler.VocabularyAccessForUser, middleware.Auth(h.updateAccessForUser))
}

func (h *Handler) changeAccess(c *gin.Context) {
	ctx := c.Request.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.Bind(&vocabAccessRq)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.addVocabulary - check body: %v", err))
		return
	}

	err = h.vocabAccessSvc.ChangeAccess(ctx, vocabAccess.Access{
		VocabID:    vocabAccessRq.VocabID,
		ID:         vocabAccessRq.ID,
		AccessEdit: vocabAccessRq.AccessEdit,
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.addVocabulary: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) addAccessForUser(c *gin.Context) {
	ctx := c.Request.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.Bind(&vocabAccessRq)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.addAccessForUser - check body: %v", err))
		return
	}

	err = h.vocabAccessSvc.AddAccessForUser(ctx, vocabAccess.Access{
		VocabID:    vocabAccessRq.VocabID,
		UserID:     vocabAccessRq.UserID,
		AccessEdit: vocabAccessRq.AccessEdit,
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.addAccessForUser: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) removeAccessForUser(c *gin.Context) {
	ctx := c.Request.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.Bind(&vocabAccessRq)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.removeAccessForUser - check body: %v", err))
		return
	}

	err = h.vocabAccessSvc.RemoveAccessForUser(ctx, vocabAccess.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.removeAccessForUser: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) updateAccessForUser(c *gin.Context) {
	ctx := c.Request.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.Bind(&vocabAccessRq)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.updateAccessForUser - check body: %v", err))
		return
	}

	err = h.vocabAccessSvc.UpdateAccessForUser(ctx, vocabAccess.Access{
		VocabID:    vocabAccessRq.VocabID,
		UserID:     vocabAccessRq.UserID,
		AccessEdit: vocabAccessRq.AccessEdit,
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.updateAccessForUser: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
