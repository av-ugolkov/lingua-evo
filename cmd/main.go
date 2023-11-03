package main

import (
	"flag"

	"lingua-evo/internal/app"
	"lingua-evo/internal/config"
	"lingua-evo/pkg/http/static"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/server_config.yaml", "it's name of application config")
	var webPath string
	flag.StringVar(&webPath, "web_path", "./website", "it's path to static files")
	var staticFilePath string
	flag.StringVar(&staticFilePath, "static_web_path", "./", "it's path to static files")

	flag.Parse()

	static.InitStaticFiles(staticFilePath)
	cfg := config.InitConfig(configPath)
	app.ServerStart(cfg, webPath)
}
