package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	ginext "github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/google/uuid"
)

func (h *Handler) adminHandler(r *ginext.Engine) {
	r.POST("/admin/vocabulary/change-translation-lang", middleware.Auth(h.changeVocabTranslationLang))
}

func (h *Handler) changeVocabTranslationLang(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			fmt.Errorf("vocabulary.delivery.Handler.changeVocabTranslationLang: %v", err)
	}

	var data struct {
		VocabID            uuid.UUID `json:"vocab_id"`
		NewTranslationLang string    `json:"lang"`
	}
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.changeVocabTranslationLang: %v", err)
	}

	err = h.vocabSvc.ChangeVocabTranslationLang(ctx, uid, data.VocabID, data.NewTranslationLang)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("vocabulary.delivery.Handler.changeVocabTranslationLang: %v", err)
	}

	return http.StatusOK, nil, nil
}
