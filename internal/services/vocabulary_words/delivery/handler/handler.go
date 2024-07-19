package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	vocabWords "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary_words"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

const (
	paramsVocabName     string = "name"
	paramsVocabID       string = "id"
	paramsPage          string = "page"
	paramsPerPage       string = "per_page"
	paramsSearch        string = "search"
	paramsOrder         string = "order"
	paramsNativeLang    string = "native_lang"
	paramsTranslateLang string = "translate_lang"
)

type Handler struct {
	vocabWordsSvc *vocabWords.Service
}

func Create(r *gin.Engine, vocabWordsSvc *vocabWords.Service) {
	h := newHandler(vocabWordsSvc)
	h.register(r)
}

func newHandler(vocabWordsSvc *vocabWords.Service) *Handler {
	return &Handler{
		vocabWordsSvc: vocabWordsSvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.POST(handler.VocabularyCopy, middleware.Auth(h.copyVocabulary))
}

func (h *Handler) copyVocabulary(c *gin.Context) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExt.SendError(c, http.StatusUnauthorized,
			fmt.Errorf("vocabulary.delivery.Handler.copyVocabulary - get user id: %v", err))
		return
	}

	vid, err := ginExt.GetQueryUUID(c, paramsVocabID)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("vocabulary.delivery.Handler.copyVocabulary - get query [id]: %v", err))
		return
	}

	err = h.vocabWordsSvc.CopyVocab(ctx, uid, vid)
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("vocabulary.delivery.Handler.copyVocabulary: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
