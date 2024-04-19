// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package main

import (
	"log/slog"
	"os"

	"github.com/bradenhc/kolob/internal"
)

func main() {
	config, err := internal.LoadServerConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Successfully loaded configuration", "port", config.Port, "data", config.DatabaseFile)
}
