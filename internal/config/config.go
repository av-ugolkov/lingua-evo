package config

import (
	"fmt"
	"sync"

	"lingua-evo/pkg/logging"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	pathConfig = "../../configs/server_config.yaml"
)

type Config struct {
	IsDebug  bool     `yaml:"is_debug"`
	JWT      JWT      `yaml:"jwt"`
	Service  Service  `yaml:"service"`
	Database Database `yaml:"database"`
	Front    Front    `yaml:"front"`
}

type JWT struct {
	Secret string `yaml:"secret" env-required:"true"`
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

type Front struct {
	Root string `yaml:"root" env-default:"./view"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig(pathConfig, instance); err != nil {
			logger.Fatalf("Fail read config: %v", err)
		}
	})
	return instance
}
