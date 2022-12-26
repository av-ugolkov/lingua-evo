package auth

import (
	"net/http"
	"os"

	"lingua-evo/internal/clients/web/entity"
	"lingua-evo/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	auth         = "/auth"
	authPagePath = entity.RootPath + "/auth/auth.html"
)

type authPage struct {
	logger *logging.Logger
}

func CreatePage(logger *logging.Logger) *authPage {
	return &authPage{
		logger: logger,
	}
}

func (p *authPage) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, auth, p.auth)
}

func (p *authPage) auth(w http.ResponseWriter, _ *http.Request) {
	file, err := os.ReadFile(authPagePath)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(file)
}
