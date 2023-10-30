package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	pathConfig = "./../configs/%s.yaml"
)

type Config struct {
	PprofDebug PprofDebug `yaml:"pprof_debug"`
	JWT        JWT        `yaml:"jwt"`
	Service    Service    `yaml:"service"`
	Database   Database   `yaml:"database"`
}

type PprofDebug struct {
	Enable bool `yaml:"enable"`
	Port   int  `yaml:"port"`
}

type JWT struct {
	Secret        string `yaml:"secret" env-required:"true"`
	ExpireAccess  int    `yaml:"expire_access" env-default:"1800"`
	ExpireRefresh int    `yaml:"expire_refresh" env-default:"2592000"`
}

type Service struct {
	Type string `yaml:"type" env-default:"tcp"`
	Port string `yaml:"port" env-default:"8080"`
}

type Database struct {
	NameDB   string `yaml:"name_db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func (db *Database) GetConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", db.User, db.Password, db.Host, db.Port, db.NameDB)
}

var instance *Config

func InitConfig(config string) *Config {
	slog.Info("read application config")
	instance = &Config{}
	fullPathConfig := fmt.Sprintf(pathConfig, config)
	if err := cleanenv.ReadConfig(fullPathConfig, instance); err != nil {
		slog.Error(fmt.Errorf("Fail read config: %v", err).Error())
		return nil
	}
	return instance
}

func GetConfig() *Config {
	return instance
}
