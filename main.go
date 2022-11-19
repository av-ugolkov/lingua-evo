package main

import (
	"log"

	tgClient "lingua-evo/clients/telegram"
	"lingua-evo/consumer/event_consumer"
	"lingua-evo/events/commands/messengers/telegram"
	"lingua-evo/storage/dictionary"
)

const (
	tgToken     = "5762950198:AAHRVBXPgAgrbSv-fUcXeAwbwDysiTXcMtY"
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	tg := tgClient.New(tgBotHost, tgToken)

	eventProcessor := telegram.New(tg, dictionary.New(storagePath))

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
