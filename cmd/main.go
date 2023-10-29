package main

import (
	"flag"

	"lingua-evo/internal/app"
	"lingua-evo/internal/config"
)

func main() {
	var nameConfig string
	flag.StringVar(&nameConfig, "config", "server_config", "it's name of application config")

	flag.Parse()

	cfg := config.InitConfig(nameConfig)
	app.ServerStart(cfg)
}
