package main

import (
	"os"

	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/config"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/di"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/di/flags"
)

func main() {
	flags := flags.ParseFlags()
	cfg := config.MustLoadConfig(*flags.ConfigPath)

	app := di.MustConfigureApp(flags, cfg)

	exitCode := app.Serve(os.Stdin, os.Stdout)

	os.Exit(exitCode)
}
