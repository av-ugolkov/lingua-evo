package api

import (
	"lingua-evo/internal/api/account"
	"lingua-evo/internal/api/add_word"
	"lingua-evo/internal/api/auth"
	"lingua-evo/internal/service"
	"lingua-evo/pkg/logging"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type api struct {
	logger *logging.Logger
	lingua *service.Lingua
}

func CreateApi(logger *logging.Logger, lingua *service.Lingua) *api {
	return &api{
		logger: logger,
		lingua: lingua,
	}
}

func (a *api) RegisterApi(router *httprouter.Router) {
	router.ServeFiles("/view/*filepath", http.Dir("./view/"))

	a.logger.Info("register auth api")
	authHandler := auth.NewHandler(a.logger)
	authHandler.Register(router)

	a.logger.Info("register account page")
	accountPage := account.NewHandler(a.logger)
	accountPage.Register(router)

	a.logger.Info("register add word page")
	addWordPage := add_word.NewHandler(a.logger, a.lingua)
	addWordPage.Register(router)

}
