package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env    string `yaml:"env" env:"ENV" env-default:"local" env-description:"app environment"`
	Minio  `yaml:"minio" env-description:"minio environment"`
	Server `yaml:"server" env-description:"server environment"`
}

type Minio struct {
	Host            string `yaml:"host" env-required:"true"`
	Port            string `yaml:"port" env-required:"true"`
	AccessKey       string `yaml:"access_key" env-required:"true"`
	SecretAccessKey string `yaml:"secret_access_key" env-required:"true"`
	SSLMode         *bool  `yaml:"ssl_mode" env-required:"true"`
	BaseBucket      string `yaml:"base_bucket" env-required:"true"`
	BasePath        string `yaml:"base_path" env-required:"true"`
}

type Server struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func Load() (*Config, error) {
	configPath := "./config/local.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot read config, path: %v, err: %w", configPath, err)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, fmt.Errorf("cannot read config err: %w", err)
	}

	return &config, nil
}
