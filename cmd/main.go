package main

import (
	"flag"

	"github.com/av-ugolkov/lingua-evo/internal/app"
	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/server_config.yaml", "it's name of application config")

	var emailPsw string
	flag.StringVar(&emailPsw, "epsw", runtime.EmptyString, "email password for newsletter")

	var jwtSecret string
	flag.StringVar(&jwtSecret, "jwts", runtime.EmptyString, "solt for jwt tokens")

	var pgPsw string
	flag.StringVar(&pgPsw, "pg_psw", runtime.EmptyString, "password for postgres db")

	var redisPsw string
	flag.StringVar(&redisPsw, "redis_psw", runtime.EmptyString, "password for redis db")

	var googleClientID string
	flag.StringVar(&googleClientID, "google_client_id", runtime.EmptyString, "google client id")

	flag.Parse()

	if jwtSecret == runtime.EmptyString ||
		pgPsw == runtime.EmptyString ||
		redisPsw == runtime.EmptyString ||
		googleClientID == runtime.EmptyString {
		panic("empty jwts, pg_psw or redis_psw")
	}

	cfg := config.InitConfig(configPath)
	config.SetEmailPassword(emailPsw)
	config.SetJWTSecret(jwtSecret)
	config.SetDBPassword(pgPsw)
	config.SetRedisPassword(redisPsw)
	config.SetGoogleClientID(googleClientID)

	app.ServerStart(cfg)
}
