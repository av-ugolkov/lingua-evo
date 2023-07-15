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

	a.logger.Info("create index")
	index.Create(a.logger, a.lingua, router)

	a.logger.Info("create auth api")
	auth.Create(a.logger, a.lingua, router)

	a.logger.Info("create account page")
	account.Create(a.logger, a.lingua, router)

	a.logger.Info("create add word page")
	add_word.Create(a.logger, a.lingua, router)
}
