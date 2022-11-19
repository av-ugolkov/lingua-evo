package main

import (
	"LinguaEvo/events/commands/messengers"
	"LinguaEvo/storage/dictionary"
	"log"

	tgClient "LinguaEvo/clients/telegram"
	"LinguaEvo/consumer/event_consumer"
)

const (
	tgToken     = "5762950198:AAHRVBXPgAgrbSv-fUcXeAwbwDysiTXcMtY"
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	tgClient := tgClient.New(tgBotHost, tgToken)

	eventProcessor := messengers.New(
		tgClient,
		dictionary.New(storagePath),
	)

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
