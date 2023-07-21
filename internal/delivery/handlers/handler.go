package handlers

import (
	"net/http"

	"lingua-evo/internal/delivery/handlers/account"
	"lingua-evo/internal/delivery/handlers/add_word"
	"lingua-evo/internal/delivery/handlers/index"
	"lingua-evo/internal/delivery/handlers/sign_in"
	"lingua-evo/internal/delivery/handlers/sign_up"
	"lingua-evo/internal/service"
	"lingua-evo/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	filePath = "/web/*filepath"
	rootPath = "./web"
)

func RegisterHandlers(logger *logging.Logger, router *httprouter.Router, lingua *service.Lingua) {
	router.ServeFiles(filePath, http.Dir(rootPath))

	logger.Info("create index")
	index.Create(logger, lingua, router)

	logger.Info("create sign_in api")
	sign_in.Create(logger, lingua, router)

	logger.Info("create sign_up api")
	sign_up.Create(logger, lingua, router)

	logger.Info("create account page")
	account.Create(logger, lingua, router)

	logger.Info("create add word page")
	add_word.Create(logger, lingua, router)
}
