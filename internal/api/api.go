package api

import (
	"lingua-evo/internal/api/auth"
	"lingua-evo/internal/delivery/repository"
	"lingua-evo/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

type api struct {
	logger  *logging.Logger
	storage repository.Storage
}

func CreateApi(logger *logging.Logger, storage repository.Storage) *api {
	return &api{
		logger:  logger,
		storage: storage,
	}
}

func (a *api) RegisterApi(router *httprouter.Router) {
	a.logger.Info("register auth api")
	authHandler := auth.NewHandler(a.logger)
	authHandler.Register(router)
}
