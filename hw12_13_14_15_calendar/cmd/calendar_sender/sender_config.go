package main

import (
	"fmt"
	"os"

	yml "gopkg.in/yaml.v3"
)

type Config struct {
	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"logger"`
	Rabbit struct {
		URL      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"rabbit"`
	EventQueue struct {
		Name     string `yaml:"name"`
		Exchange string `yaml:"exchange"`
	} `yaml:"event-queue"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	return &config, nil
}
