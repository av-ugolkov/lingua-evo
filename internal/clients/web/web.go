package web

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"lingua-evo/internal/clients/web/pages/auth"
	"lingua-evo/internal/config"
	"lingua-evo/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

type Web struct {
	logger logging.Logger
}

func CreateWeb() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	var server *http.Server
	var listener net.Listener

	router := httprouter.New()
	registerHandlers(router, logger)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s: %s", cfg.Listen.BindIP, cfg.Listen.Port))
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("listener: %v", listener.Addr())

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
	logger.Print("register auth")
	authHandler := auth.NewHandler(logger)
	authHandler.Register(router)
}
