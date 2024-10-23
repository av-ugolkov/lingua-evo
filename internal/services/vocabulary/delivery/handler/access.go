package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/gin-gonic/gin"
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

func (h *Handler) getAccessForUser(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getAccessForUser: %v", err)
	}

	vid, err := c.GetQueryUUID(paramsID)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getAccessForUser: %v", err)
	}

	accessVocab, err := h.vocabSvc.GetAccessForUser(ctx, uid, vid)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.getAccessForUser: %v", err)
	}

	return http.StatusOK, gin.H{
		"access": accessVocab,
	}, nil
}

func (h *Handler) addAccessForUser(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.Bind(&vocabAccessRq)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("vocabulary.delivery.Handler.addAccessForUser: %v", err)
	}

	err = h.vocabSvc.AddAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
		Status:  vocabAccessRq.Status,
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.addAccessForUser: %v", err)
	}

	return http.StatusOK, gin.H{}, nil
}

func (h *Handler) removeAccessForUser(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.Bind(&vocabAccessRq)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("vocabulary.delivery.Handler.removeAccessForUser: %v", err)
	}

	err = h.vocabSvc.RemoveAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.removeAccessForUser: %v", err)
	}

	return http.StatusOK, gin.H{}, nil
}

func (h *Handler) updateAccessForUser(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var vocabAccessRq VocabularyAccessRq
	err := c.Bind(&vocabAccessRq)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("vocabulary.delivery.Handler.updateAccessForUser: %v", err)
	}

	err = h.vocabSvc.UpdateAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
		Status:  vocabAccessRq.Status,
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.updateAccessForUser: %v", err)
	}

	return http.StatusOK, gin.H{}, nil
}
