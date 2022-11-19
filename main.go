package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"

	tgClient "lingua-evo/clients/telegram"
	"lingua-evo/consumer/event_consumer"
	"lingua-evo/events/commands/messengers/telegram"
	"lingua-evo/storage"
	"lingua-evo/storage/database"
	"lingua-evo/storage/dictionary"
)

const (
	tgToken      = "5762950198:AAHRVBXPgAgrbSv-fUcXeAwbwDysiTXcMtY"
	tgBotHost    = "api.telegram.org"
	storagePath  = "storage"
	dbConnection = "postgres://postgres:5623@localhost:5432/postgres"
	batchSize    = 100
)

func main() {
	tg := tgClient.New(tgBotHost, tgToken)

	var repository storage.Storage
	pool, err := pgxpool.Connect(context.Background(), dbConnection)
	if err != nil {
		log.Printf("can't create pg pool: %v", err)
		repository = dictionary.New(storagePath)
	} else {
		log.Printf("create pg pool: %v", pool.Config().ConnConfig.Database)
		repository = database.New(pool)
	}

	eventProcessor := telegram.New(tg, repository)

	log.Print("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
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
