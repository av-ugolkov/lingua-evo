package main

import (
	"flag"
	"log/slog"

	"github.com/av-ugolkov/lingua-evo/internal/app"
	"github.com/av-ugolkov/lingua-evo/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/server_config.yaml", "it's name of application config")

	var emailPsw string
	flag.StringVar(&emailPsw, "epsw", "", "email password for newsletter")

	var jwtSecret string
	flag.StringVar(&jwtSecret, "jwts", "", "")

	var pgPsw string
	flag.StringVar(&pgPsw, "pg_psw", "", "")

	var redisPsw string
	flag.StringVar(&redisPsw, "redis_psw", "", "")

	flag.Parse()

	if jwtSecret == "" || pgPsw == "" || redisPsw == "" {
		slog.Error("empty jwts, pg_psw or redis_psw")
		return
	}

	cfg := config.InitConfig(configPath)
	config.SetEmailPassword(emailPsw)
	config.SetJWTSecret(jwtSecret)
	config.SetDBPassword(pgPsw)
	config.SetRedisPassword(redisPsw)

	app.ServerStart(cfg)
}
