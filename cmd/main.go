package main

import (
	"os"

	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/config"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/di"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/di/flags"
	exitcodes "github.com/SmokingElk/MWS-2025-Autumn-Intership/pkg/exit-codes"
)

func main() {
	flags := flags.ParseFlags()
	cfg := config.MustLoadConfig(*flags.ConfigPath)

	app := di.MustConfigureApp(flags, cfg)

	err := app.Serve()

	if err == nil {
		os.Exit(exitcodes.OK)
	} else {
		os.Exit(exitcodes.UnknownError)
	}
}
