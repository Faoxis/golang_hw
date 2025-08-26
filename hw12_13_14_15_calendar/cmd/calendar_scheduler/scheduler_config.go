package main

import (
	"fmt"
	"os"

	yml "gopkg.in/yaml.v3"
)

type SchedulerConfig struct {
	Logger     Logger
	Storage    Storage
	Rabbit     Rabbit
	EventQueue EventQueue `yaml:"event-queue"`
	Scheduler  SchedulerSettings
}

type EventQueue struct {
	Name     string
	Exchange string
}

type SchedulerSettings struct {
	CheckInterval string `yaml:"check-interval"`
}

type Rabbit struct {
	Url      string
	Username string
	Password string
}

type Storage struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string `yaml:"sslmode"`
}

type Logger struct {
	Level string
}

func LoadConfig(path string) (*SchedulerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := yml.NewDecoder(file)
	var cfg SchedulerConfig
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	return &cfg, nil
}

func (storage *Storage) GetPostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		storage.Host,
		storage.Port,
		storage.User,
		storage.Password,
		storage.Database,
		storage.SSLMode,
	)
}
