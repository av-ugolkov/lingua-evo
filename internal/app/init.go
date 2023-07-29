package app

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/delivery/handlers"
	"lingua-evo/internal/delivery/repository"
	"lingua-evo/internal/services"

	"lingua-evo/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

func ServerStart(logger *logging.Logger, cfg *config.Config) {
	db, err := repository.NewDB(cfg.Database.GetConnStr())
	if err != nil {
		logger.Fatalf("can't create pg pool: %v", err)
	}

	database := repository.NewDatabase(db)
	lingua := services.NewLinguaService(logger, database)

	router := httprouter.New()

	handlers.RegisterHandlers(logger, router, lingua)

	address := fmt.Sprintf(":%s", cfg.Service.Port)

	listener, err := net.Listen(cfg.Service.Type, address)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("web address: %v", listener.Addr())

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
