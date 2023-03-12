package add_word

import (
	"github.com/julienschmidt/httprouter"
	"lingua-evo/pkg/logging"
	"net/http"
)

const (
	addWordURL = "/add_word"
)

type Handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, addWordURL, h.addWord)
}
