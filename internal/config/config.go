package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	CLIConfig    `yaml:"cli" env-required:"true"`
	GithubConfig `yaml:"github" env-required:"true"`
}

type CLIConfig struct {
	TimeoutSeconds int `yaml:"timeout_seconds" env-required:"true"`
}

type GithubConfig struct {
	AuthToken string `yaml:"auth_token"`
}

func MustLoadConfig(configPath string) *Config {
	if configPath == "" {
		log.Fatal("config-path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	return &cfg
}
