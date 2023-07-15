package api

import (
	"fmt"
	"net/http"

	"lingua-evo/internal/delivery/api/account"
	"lingua-evo/internal/delivery/api/add_word"
	"lingua-evo/internal/delivery/api/auth"
	"lingua-evo/internal/delivery/api/index"
	"lingua-evo/internal/service"
	"lingua-evo/pkg/logging"
	"lingua-evo/pkg/tools/view"

	"github.com/julienschmidt/httprouter"
)

const (
	pathFolder = "view/"
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
	path := view.GetPathFolder(fmt.Sprintf("%s%s", pathFolder, "*filepath"))
	root := http.Dir(view.GetPathFile(pathFolder))
	a.logger.Debugf("%s ::: %s", root, path)
	router.ServeFiles(path, root)

	a.logger.Info("register index")
	indexHandler := index.NewHandler(a.logger, a.lingua)
	indexHandler.Register(router)

	a.logger.Info("register auth api")
	authHandler := auth.NewHandler(a.logger, a.lingua)
	authHandler.Register(router)

	a.logger.Info("register account page")
	accountPage := account.NewHandler(a.logger)
	accountPage.Register(router)

	a.logger.Info("register add word page")
	addWordPage := add_word.NewHandler(a.logger, a.lingua)
	addWordPage.Register(router)
}
