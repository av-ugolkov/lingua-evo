package auth

import (
	"net/http"
	"os"

	"lingua-evo/internal/clients/web/entity"
	"lingua-evo/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	auth = "/auth"
)

type authPage struct {
	logger *logging.Logger
}

func CreatePage(logger *logging.Logger) *authPage {
	return &authPage{
		logger: logger,
	}
}

func (a *authPage) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, auth, a.auth)
}

func (a *authPage) auth(w http.ResponseWriter, _ *http.Request) {
	file, err := os.ReadFile("./pages/auth/auth.html")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

type Page struct {
	Title string
	Body  []byte
}
