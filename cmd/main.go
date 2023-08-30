package main

import (
	"lingua-evo/internal/app"
	"lingua-evo/internal/config"
	"log/slog"
)

func main() {
	cfg := config.GetConfig()
	slog.Info("config initializing")

	app.ServerStart(cfg)
}
