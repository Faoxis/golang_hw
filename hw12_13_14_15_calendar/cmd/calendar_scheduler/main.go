package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	fmt.Println("Starting calendar scheduler")

	// Загружаем конфигурацию
	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error: loading config from file %d", err)
	}

	// Инициализируем логгер
	logg := logger.New(config.Logger.Level)

	_, err = initDatabaseStorage(config, logg)
	if err != nil {
		log.Fatalf("Error: initializing database storage %d", err)
	}

}

func initDatabaseStorage(config *SchedulerConfig, logg app.Logger) (app.Storage, error) {
	migrationsPath, err := filepath.Abs("./migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to get migrations path: %w", err)
	}

	// Выполняем миграции
	if err := sqlstorage.RunMigrations(config.Storage.GetPostgresDSN(), migrationsPath); err != nil {
		return nil, fmt.Errorf("migrations failed: %w", err)
	}

	// Создаем хранилище
	storage, err := sqlstorage.New(config.Storage.GetPostgresDSN(), logg)
	if err != nil {
		return nil, fmt.Errorf("sql storage failed: %w", err)
	}

	return storage, nil
}
