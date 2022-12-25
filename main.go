package main

import (
	"context"

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

	var repository storage.Storage
	pool, err := pgxpool.Connect(context.Background(), dbConnection)
	if err != nil {
		logger.Fatalf("can't create pg pool: %v", err)
	}
	logger.Printf("create pg pool: %v", pool.Config().ConnConfig.Database)
	repository = database.New(pool)
	logger.Printf("repository: %s", repository)

	//eventProcessor := telegram.New(tg, repository)

	//log.Print("service started")

	//consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)
	//if err := consumer.Start(); err != nil {
	//	log.Fatal("service is stopped", err)
	//}

	web.CreateWeb(logger, cfg)
}

/*func mustToken() string {
	token := flag.String("tg-token-bot", "", "token for access to telegram bot")
	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}*/
