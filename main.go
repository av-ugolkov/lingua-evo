package main

import (
	"lingua-evo/internal/app"
	"lingua-evo/internal/config"

	"lingua-evo/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	cfg := config.GetConfig()
	logger.Println("config initializing")

	app.ServerStart(logger, cfg)
}
