// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package server

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	Port            int
	DatabaseFile    string
	ShutdownTimeout time.Duration
}

func LoadConfig() (Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		var s Config
		return s, err
	}

	s := Config{
		Port:            24000,
		DatabaseFile:    path.Join(cwd, "kolob.db"),
		ShutdownTimeout: 10 * time.Second,
	}
	s.loadEnvironment()
	s.loadArgs()
	return s, nil
}

func (s *Config) loadEnvironment() error {
	if val := os.Getenv("KOLOB_PORT"); val != "" {
		port, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("failed to parse KOLOB_PORT: %v", err)
		}
		s.Port = port
	}

	if val := os.Getenv("KOLOB_DATA"); val != "" {
		s.DatabaseFile = val
	}

	if val := os.Getenv("KOLOB_SHUTDOWN_TIMEOUT"); val != "" {
		d, err := time.ParseDuration(val)
		if err != nil {
			return fmt.Errorf("failed to parse KOLOB_SHUTDOWN_TIMEOUT: %v", err)
		}
		s.ShutdownTimeout = d
	}

	return nil
}

func (s *Config) loadArgs() error {
	port := flag.Int("port", 0, "The port to run the HTTP server on.")
	data := flag.String("data", "", "The path to the database file where data is stored.")

	flag.Usage = func() {
		println := func(format string, a ...any) {
			fmt.Fprintf(flag.CommandLine.Output(), format, a...)
			fmt.Fprint(flag.CommandLine.Output(), "\n")
		}

		println("")
		println("usage:  %s [options...]", filepath.Base(os.Args[0]))
		println("")
		println("Kolob is a lightweight and secure collaboration server.")
		println("")
		flag.PrintDefaults()
		println("")
	}

	flag.Parse()

	if *port != 0 {
		s.Port = *port
	}
	if *data != "" {
		s.DatabaseFile = *data
	}
	return nil
}
