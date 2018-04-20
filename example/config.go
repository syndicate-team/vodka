package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/niklucky/go-lib"
	"github.com/niklucky/vodka"
	"github.com/niklucky/vodka/adapters"
)

/*
Config - service configuration
*/
type Config struct {
	Version    string           `json:"version"`
	HTTPServer vodka.HTTPConfig `json:"http_server"`
	Postgres   adapters.Config  `json:"postgres"`
	Debug      bool
}

/*
NewConfig - config constructors
*/
func NewConfig(configFileName string) (Config, error) {
	config := Config{}

	fileData, err := lib.ReadFile(configFileName)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return config, err
	}
	if os.Getenv("DEBUG") != "" {
		config.Debug = true
		log.Println("Running in debug mode")
	}
	return config, nil
}
