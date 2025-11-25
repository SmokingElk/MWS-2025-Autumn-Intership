package main

import (
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/config"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/di"
)

func main() {
	cfg := config.MustLoadConfig()

	di.MustConfigureApp(cfg)
}
