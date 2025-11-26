package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GithubConfig `yaml:"github" env-required:"true"`
}

type GithubConfig struct {
	AuthToken string `yaml:"auth_token"`
}

func MustLoadConfig(configPath string) *Config {
	var cfg Config

	if configPath == "" {
		return &cfg
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	return &cfg
}
