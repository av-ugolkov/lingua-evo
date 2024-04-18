package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/av-ugolkov/lingua-evo/internal/app"
	"github.com/av-ugolkov/lingua-evo/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/server_config.yaml", "it's name of application config")

	flag.Parse()

	slog.Info(fmt.Sprintf("config path: %s", configPath))
	cfg := config.InitConfig(configPath)
	app.ServerStart(cfg)
}
