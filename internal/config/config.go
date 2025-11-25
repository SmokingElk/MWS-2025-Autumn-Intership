package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	GithubConfig `yaml:"github" env-required:"true"`
}

type GithubConfig struct {
	AuthToken string `yaml:"auth_token"`
}

func MustLoadConfig() *Config {
	if os.Getenv("CONFIG_PATH") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("failed to load .env")
		}
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
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
