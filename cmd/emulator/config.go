package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jonathangjertsen/toyboy/model"
)

type Config struct {
	Location string
	Model    model.Config
}

func LoadConfig(location string) (Config, error) {
	var config Config
	jsonData, err := os.ReadFile(location)
	if err != nil {
		return config, fmt.Errorf("failed to load config file in %s", location)
	}
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return config, fmt.Errorf("config in %s is corrupted", location)
	}
	config.Location = location
	return config, nil
}

func (conf *Config) Save() {
	jsonData, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		fmt.Printf("marshalling config: %v", err)
		return
	}
	if err := os.WriteFile(conf.Location, jsonData, 0o666); err != nil {
		fmt.Printf("failed to save config file: %v", err)
	}
}
