package config

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	pathConfig = "./../configs/server_config.yaml"
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
	ExpireAccess  int    `yaml:"expire_access" env-default:"60*60"`
	ExpireRefresh int    `yaml:"expire_refresh" env-default:"60*60*24*30"`
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
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		slog.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig(pathConfig, instance); err != nil {
			slog.Error(fmt.Errorf("Fail read config: %v", err).Error())
			return
		}
	})
	return instance
}
