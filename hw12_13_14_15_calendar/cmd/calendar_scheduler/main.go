package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	fmt.Println("Starting calendar scheduler")

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error: loading config from file %v", err)
	}

	logg := logger.New(config.Logger.Level)

	notificationStorage, err := initNotificationStorage(config, logg)
	if err != nil {
		log.Fatalf("Error: initializing notification storage %v", err)
	}
	defer notificationStorage.Close()

	queue, err := rabbit.NewRabbitQueue[storage.Notification](
		config.Rabbit.Url,
		config.Rabbit.Username,
		config.Rabbit.Password,
		logg,
	)
	if err != nil {
		log.Fatalf("Error: creating rabbit queue %v", err)
	}
	defer queue.Close()

	scheduler := scheduler.NewScheduler(logg, notificationStorage, queue, config.EventQueue.Name, config.EventQueue.Exchange, config.Scheduler.CheckInterval)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go scheduler.Start(ctx)

	<-sigChan
	fmt.Println("Scheduler stopped")
}

func initNotificationStorage(config *SchedulerConfig, logg app.Logger) (scheduler.NotificationStorage, error) {
	migrationsPath, err := filepath.Abs("./migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to get migrations path: %w", err)
	}

	if err := sqlstorage.RunMigrations(config.Storage.GetPostgresDSN(), migrationsPath); err != nil {
		return nil, fmt.Errorf("migrations failed: %w", err)
	}

	notificationStorage, err := scheduler.NewSQLNotificationStorage(config.Storage.GetPostgresDSN(), logg)
	if err != nil {
		return nil, fmt.Errorf("sql notification storage failed: %w", err)
	}

	return notificationStorage, nil
}
