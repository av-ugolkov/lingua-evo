package add_word

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"lingua-evo/internal/service"
	"lingua-evo/pkg/logging"
)

const (
	addWordURL = "/add_word"
)

type Handler struct {
	logger *logging.Logger
	lingua *service.Lingua
}

func NewHandler(logger *logging.Logger, lingua *service.Lingua) *Handler {
	return &Handler{
		logger: logger,
		lingua: lingua,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, addWordURL, h.getAddWord)
	router.HandlerFunc(http.MethodPost, addWordURL, h.postAddWord)
}
