package main

import (
	"context"
	"flag"
	sqlstorage "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error: loading config from file %d", err)
	}
	logg := logger.New(config.Logger.Level)

	var storage app.Storage
	switch config.Storage.Type {
	case "memory":
		storage = memorystorage.New(logg)
	case "database":
		migrationsPath, err := filepath.Abs("./migrations")
		if err != nil {

		}
		err = sqlstorage.RunMigrations(
			config.Storage.GetPostgresDSN(),
			migrationsPath,
		)
		if err != nil {
			log.Fatalf("migrations failed: %v", err)
		}
		storage, err = sqlstorage.New(config.Storage.GetPostgresDSN(), logg)
		if err != nil {
			log.Fatalf("sql storage failed: %v", err)
		}
	}
	defer storage.Close()
	calendar := app.New(logg, storage)

	host := config.Server.Host
	port := config.Server.Port
	server := internalhttp.NewServer(logg, host, port, calendar)
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
