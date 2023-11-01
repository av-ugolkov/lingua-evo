package main

import (
	"flag"
	"fmt"
	"log/slog"

	"lingua-evo/internal/app"
	"lingua-evo/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/server_config.yaml", "it's name of application config")
	var webPath string
	flag.StringVar(&webPath, "web_path", "./website", "it's path to static files")

	flag.Parse()

	slog.Info(fmt.Sprintf("configPath: %s", configPath))
	slog.Info(fmt.Sprintf("webPath: %s", webPath))

	cfg := config.InitConfig(configPath)
	app.ServerStart(cfg, webPath)
}
