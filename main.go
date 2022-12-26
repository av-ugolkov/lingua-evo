package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"lingua-evo/internal/api"
	"net"
	"net/http"
	"time"

	"lingua-evo/internal/clients/web"
	"lingua-evo/internal/config"
	"lingua-evo/pkg/logging"
	"lingua-evo/pkg/storage"
	"lingua-evo/pkg/storage/database"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	tgToken      = "5762950198:AAHRVBXPgAgrbSv-fUcXeAwbwDysiTXcMtY"
	tgBotHost    = "api.telegram.org"
	dbConnection = "postgres://postgres:5623@localhost:5432/postgres"
	batchSize    = 100
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	logger.Println("config initializing")
	cfg := config.GetConfig()

	//tg := tgClient.New(tgBotHost, tgToken)

	pool, err := pgxpool.Connect(context.Background(), dbConnection)
	if err != nil {
		logger.Fatalf("can't create pg pool: %v", err)
	}
	logger.Printf("create pg pool: %v", pool.Config().ConnConfig.Database)

	var repository storage.Storage
	repository = database.New(pool)

	//eventProcessor := telegram.New(tg, repository)

	//log.Print("service started")

	//consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)
	//if err := consumer.Start(); err != nil {
	//	log.Fatal("service is stopped", err)
	//}
	router := httprouter.New()

	api := api.CreateApi(logger, repository)
	api.RegisterApi(router)
	web := web.CreateWeb(logger)
	web.Register(router)

	createServer(router, logger, cfg)
}

func createServer(router *httprouter.Router, logger *logging.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

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

/*func mustToken() string {
	token := flag.String("tg-token-bot", "", "token for access to telegram bot")
	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}*/
