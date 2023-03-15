package config

import (
	"fmt"
	"sync"

	"lingua-evo/pkg/logging"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug  *bool    `yaml:"is_debug"`
	JWT      JWT      `yaml:"jwt"`
	Listen   Listen   `yaml:"listen"`
	Database Database `yaml:"database"`
}

type JWT struct {
	Secret string `yaml:"secret" env-required:"true"`
}

type Listen struct {
	Type   string `yaml:"type" env-default:"port"`
	BindIP string `yaml:"bind_ip" env-default:"localhost"`
	Port   string `yaml:"port" env-default:"8080"`
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
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("configs/server_config.yaml", instance); err != nil {
			_, err := cleanenv.GetDescription(instance, nil)
			if err != nil {
				logger.Error(err)
			}
			logger.Fatal(err)
		}
	})
	return instance
}
