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

	//var emailPassword string
	//flag.StringVar(&emailPassword, "e-psw", "", "email password for newsletter")

	flag.Parse()

	slog.Info(fmt.Sprintf("config path: %s", configPath))
	cfg := config.InitConfig(configPath)

	//config.SetEmailPassword(emailPassword)

	app.ServerStart(cfg)
}
