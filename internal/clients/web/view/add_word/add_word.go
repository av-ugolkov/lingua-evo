package add_word

import (
	"github.com/julienschmidt/httprouter"
	"lingua-evo/internal/clients/web/entity"
	"lingua-evo/pkg/logging"
	"net/http"
	"os"
)

const (
	account         = "/account"
	addWordPagePath = entity.RootPath + "/account/account.html"
)

type addWordPage struct {
	logger *logging.Logger
}

func CreatePage(logger *logging.Logger) *addWordPage {
	return &addWordPage{
		logger: logger,
	}
}

func (p *addWordPage) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, account, p.account)
}

func (p *addWordPage) account(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile(addWordPagePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(file)
}
