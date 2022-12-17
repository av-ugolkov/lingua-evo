package config

import (
	"sync"

	"lingua-evo/pkg/logging"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug    *bool      `yaml:"is_debug"`
	JWT        JWT        `yaml:"jwt"`
	Listen     Listen     `json:"listen"`
	WebService WebService `yaml:"web_service" env-required:"true"`
}

type JWT struct {
	Secret string `yaml:"secret" env-required:"true"`
}

type Listen struct {
	Type   string `yaml:"type" env-default:"port"`
	BindIP string `yaml:"bind_ip" env-default:"localhost"`
	Port   string `yaml:"port" env-default:"8080"`
}

type WebService struct {
	URL string `yaml:"url" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("configs/dev.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
