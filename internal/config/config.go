package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PprofDebug PprofDebug `yaml:"pprof_debug"`
	Logger     Logger     `yaml:"logger"`
	SSL        SSL        `yaml:"ssl"`
	JWT        JWT        `yaml:"jwt"`
	Service    Service    `yaml:"service"`
	DbSQL      DbSQL      `yaml:"postgres"`
	DbRedis    DbRedis    `yaml:"redis"`
	Email      Email      `yaml:"email"`
	Kafka      Kafka      `yaml:"kafka"`
	Google     Google     `yaml:"google"`
	AES        Aes        `yaml:"aes"`
}

var instance *Config

func InitConfig(pathConfig string) *Config {
	slog.Info("read application config")
	instance = &Config{}
	if err := cleanenv.ReadConfig(pathConfig, instance); err != nil {
		slog.Error(fmt.Errorf("fail read config: %v", err).Error())
		return nil
	}

	return instance
}

func SetEmailPassword(emailPsw string) {
	instance.Email.Password = emailPsw
}

func SetJWTSecret(secret string) {
	instance.JWT.Secret = secret
}

func SetDBPassword(psw string) {
	instance.DbSQL.Password = psw
}

func SetRedisPassword(psw string) {
	instance.DbRedis.Password = psw
}

func GetConfig() *Config {
	return instance
}
