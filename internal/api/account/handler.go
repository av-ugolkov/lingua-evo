package account

import (
	"lingua-evo/pkg/logging"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	account = "/account"
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
	router.HandlerFunc(http.MethodGet, account, h.account)
}
