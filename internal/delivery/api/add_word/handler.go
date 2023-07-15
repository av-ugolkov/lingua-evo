package add_word

import (
	"net/http"

	"lingua-evo/internal/service"
	"lingua-evo/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	addWordURL = "/add_word"
)

type Handler struct {
	logger *logging.Logger
	lingua *service.Lingua
}

func Create(log *logging.Logger, ling *service.Lingua, r *httprouter.Router) {
	handler := newHandler(log, ling)
	handler.register(r)
}

func newHandler(logger *logging.Logger, lingua *service.Lingua) *Handler {
	return &Handler{
		logger: logger,
		lingua: lingua,
	}
}

func (h *Handler) register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, addWordURL, h.getAddWord)

	router.HandlerFunc(http.MethodPost, addWordURL, h.postAddWord)
}
