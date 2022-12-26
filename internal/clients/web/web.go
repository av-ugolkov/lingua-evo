package web

import (
	"net/http"

	"lingua-evo/internal/clients/web/view/account"
	"lingua-evo/internal/clients/web/view/add_word"
	authPage "lingua-evo/internal/clients/web/view/auth"
	"lingua-evo/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

type web struct {
	logger *logging.Logger
}

func CreateWeb(logger *logging.Logger) *web {
	return &web{
		logger: logger,
	}
}

func (w *web) Register(router *httprouter.Router) {
	router.ServeFiles("/pages/*filepath", http.Dir("./pages/"))

	w.logger.Info("register auth page")
	authPage := authPage.CreatePage(w.logger)
	authPage.Register(router)

	w.logger.Info("register account page")
	accountPage := account.CreatePage(w.logger)
	accountPage.Register(router)

	w.logger.Info("register add word page")
	addWordPage := add_word.CreatePage(w.logger)
	addWordPage.Register(router)
}
