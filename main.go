package main

import (
	"errors"
	"fmt"
	"lingua-evo/internal/service"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"lingua-evo/internal/api"
	"lingua-evo/internal/config"
	"lingua-evo/internal/delivery/repository"
	"lingua-evo/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	cfg := config.GetConfig()
	logger.Println("config initializing")

	//tg := tgClient.New(tgBotHost, tgToken)

	db, err := repository.NewDB(cfg.Database.GetConnStr())
	if err != nil {
		logger.Fatalf("can't create pg pool: %v", err)
	}

	//eventProcessor := telegram.New(tg, repository)

	//log.Print("service started")

	//consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)
	//if err := consumer.Start(); err != nil {
	//	log.Fatal("service is stopped", err)
	//}

	database := repository.NewDatabase(db)
	wordsService := service.NewWordsService(database)

	lingua := service.NewLinguaService(wordsService)

	router := httprouter.New()

	api := api.CreateApi(logger, lingua)
	api.RegisterApi(router)

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
