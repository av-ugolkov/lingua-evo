package main

import (
	"log/slog"

	"lingua-evo/internal/app"
	"lingua-evo/internal/config"
)

func main() {
	cfg := config.GetConfig()
	slog.Info("config initializing")

	app.ServerStart(cfg)
}
