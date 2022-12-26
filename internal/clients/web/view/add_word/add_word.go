package add_word

import (
	"net/http"
	"os"

	"lingua-evo/internal/clients/web/entity"
	"lingua-evo/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	addWord         = "/add-word"
	addWordPagePath = entity.RootPath + "/add_word/add_word.html"
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
	router.HandlerFunc(http.MethodGet, addWord, p.account)
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
