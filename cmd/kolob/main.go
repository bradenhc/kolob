// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package main

import (
	"log/slog"
	"os"

	"github.com/bradenhc/kolob/internal/server"
)

func main() {
	config, err := server.LoadConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info(
		"Successfully loaded configuration",
		slog.Group("config",
			"port", config.Port,
			"data", config.DatabaseFile,
		),
	)
	server, err := server.NewServer(config)
	if err != nil {
		slog.Error("failed to start server", "err", err.Error())
	}

	server.Start()
}
