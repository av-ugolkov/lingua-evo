package main

import (
	"flag"

	"github.com/av-ugolkov/lingua-evo/internal/app"
	"github.com/av-ugolkov/lingua-evo/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/server_config.yaml", "it's name of application config")

	var emailPassword string
	flag.StringVar(&emailPassword, "epsw", "", "email password for newsletter")

	flag.Parse()

	cfg := config.InitConfig(configPath)

	config.SetEmailPassword(emailPassword)
	app.ServerStart(cfg)
}
