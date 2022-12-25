package web

import (
	"errors"
	"fmt"
	"lingua-evo/internal/clients/web/api/auth"
	authPage "lingua-evo/internal/clients/web/pages/auth"
	"lingua-evo/internal/config"
	"lingua-evo/pkg/logging"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Web struct {
	logger logging.Logger
}

func CreateWeb(logger *logging.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

	router := httprouter.New()

	router.ServeFiles("/pages/*filepath", http.Dir("./pages/"))

	registerHandlers(router, logger)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("web address: %v", listener.Addr())

	server = &http.Server{
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

func registerHandlers(router *httprouter.Router, logger *logging.Logger) {
	logger.Info("register auth api")
	authHandler := auth.NewHandler(logger)
	authHandler.Register(router)

	logger.Info("register auth page")
	authPage := authPage.CreatePage(logger)
	authPage.Register(router)
}
