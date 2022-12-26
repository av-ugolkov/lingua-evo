package account

import (
	"net/http"
	"os"

	"lingua-evo/internal/clients/web/entity"
	"lingua-evo/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	account         = "/account"
	accountPagePath = entity.RootPath + "/account/account.html"
)

type accountPage struct {
	logger *logging.Logger
}

func CreatePage(logger *logging.Logger) *accountPage {
	return &accountPage{
		logger: logger,
	}
}

func (p *accountPage) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, account, p.account)
}

func (p *accountPage) account(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile(accountPagePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(file)
}
