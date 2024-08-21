package handler

import (
	"fmt"
	"net/http"

	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
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

func (h *Handler) getAccessForUser(c *gin.Context) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("vocabulary.delivery.Handler.getAccessForUser - unauthorized: %v", err))
		return
	}

	vid, err := ginExt.GetQueryUUID(c, paramsID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.getAccessForUser - check query [vocab_id]: %v", err))
		return
	}

	access, err := h.vocabSvc.GetAccessForUser(ctx, uid, vid)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.getAccessForUser: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access": access,
	})
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

	err = h.vocabSvc.AddAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
		Status:  vocabAccessRq.Status,
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

	err = h.vocabSvc.RemoveAccessForUser(ctx, vocabulary.Access{
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

	err = h.vocabSvc.UpdateAccessForUser(ctx, vocabulary.Access{
		VocabID: vocabAccessRq.VocabID,
		UserID:  vocabAccessRq.UserID,
		Status:  vocabAccessRq.Status,
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.updateAccessForUser: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
