package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PprofDebug PprofDebug `yaml:"pprof_debug"`
	JWT        JWT        `yaml:"jwt"`
	Service    Service    `yaml:"service"`
	DbSQL      DbSQL      `yaml:"postgres"`
	DbRedis    DbRedis    `yaml:"redis"`
	Email      Email      `yaml:"email"`
	Kafka      Kafka      `yaml:"kafka"`
}

type PprofDebug struct {
	Enable bool `yaml:"enable"`
	Port   int  `yaml:"port"`
}

func (p PprofDebug) Addr() string {
	return fmt.Sprintf("localhost:%d", p.Port)
}

type JWT struct {
	Secret        string `yaml:"secret" env-required:"true"`
	ExpireAccess  int    `yaml:"expire_access" env-default:"300"`
	ExpireRefresh int    `yaml:"expire_refresh" env-default:"2592000"`
}

type Service struct {
	Type           string   `yaml:"type" env-default:"tcp"`
	Port           string   `yaml:"port" env-default:"8080"`
	AllowedOrigins []string `yaml:"allowed_origins" env-default:"http://localhost:5173"`
}

type DbSQL struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func (db *DbSQL) GetConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", db.User, db.Password, db.Host, db.Port, db.Name)
}

type DbRedis struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       int    `yaml:"db"`
}

type Email struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password,omitempty"`
}

type Kafka struct {
	Enable bool     `yaml:"enable"`
	Host   string   `yaml:"host"`
	Port   int      `yaml:"port"`
	Topics []string `yaml:"topics"`
}

func (k Kafka) Addr() string {
	return fmt.Sprintf("%s:%d", k.Host, k.Port)
}

var instance *Config

func InitConfig(pathConfig string) *Config {
	slog.Info("read application config")
	instance = &Config{}
	if err := cleanenv.ReadConfig(pathConfig, instance); err != nil {
		slog.Error(fmt.Errorf("Fail read config: %v", err).Error())
		return nil
	}

	return instance
}

func SetEmailPassword(emailPsw string) {
	instance.Email.Password = emailPsw
}

func GetConfig() *Config {
	return instance
}
