package internal

import (
	"flag"
	"os"
	"path"
	"strconv"
)

type ServerConfig struct {
	Port         int
	DatabaseFile string
}

func LoadServerConfig() (ServerConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		var s ServerConfig
		return s, err
	}

	s := ServerConfig{
		Port:         24000,
		DatabaseFile: path.Join(cwd, "kolob.db"),
	}
	s.loadEnvironment()
	s.loadArgs()
	return s, nil
}

func (s *ServerConfig) loadEnvironment() error {
	if val := os.Getenv("KOLOB_PORT"); val != "" {
		port, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		s.Port = port
	}
	if val := os.Getenv("KOLOB_DATA"); val != "" {
		s.DatabaseFile = val
	}
	return nil
}

func (s *ServerConfig) loadArgs() error {
	port := flag.Int("port", 0, "The port to run the HTTP server on.")
	data := flag.String("data", "", "The path to the database file where data is stored.")
	flag.Parse()

	if *port != 0 {
		s.Port = *port
	}
	if *data != "" {
		s.DatabaseFile = *data
	}
	return nil
}
