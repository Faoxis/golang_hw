package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	internalgrpc "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/server/grpc"
	sqlstorage "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage/sql"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/server"
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

	// Загружаем конфигурацию
	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error: loading config from file %d", err)
	}

	// Инициализируем логгер
	logg := logger.New(config.Logger.Level)

	// Инициализируем хранилище
	storage, err := initStorage(config, logg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()

	// Создаем приложение
	calendar := app.New(logg, storage)

	// Создаем контекст для graceful shutdown
	ctx, cancel := createShutdownContext()
	defer cancel()

	// Запускаем HTTP сервер
	httpServer := initHTTPServer(config, logg, calendar)
	go gracefulShutdown(ctx, httpServer, logg)

	// Запускаем gRPC сервер
	go startGRPCServer(config, logg, calendar, ctx)

	logg.Info("calendar is running...")

	// Запускаем HTTP сервер и ждем завершения
	if err := httpServer.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}

	logg.Info("calendar is stopped")
}

// initStorage инициализирует хранилище в зависимости от конфигурации
func initStorage(config *Config, logg app.Logger) (app.Storage, error) {
	switch config.Storage.Type {
	case "memory":
		return memorystorage.New(logg), nil
	case "database":
		return initDatabaseStorage(config, logg)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", config.Storage.Type)
	}
}

// initDatabaseStorage инициализирует базу данных
func initDatabaseStorage(config *Config, logg app.Logger) (app.Storage, error) {
	// Создаем хранилище (миграции выполняются отдельным контейнером)
	storage, err := sqlstorage.New(config.Storage.GetPostgresDSN(), logg)
	if err != nil {
		return nil, fmt.Errorf("sql storage failed: %w", err)
	}

	return storage, nil
}

// initHTTPServer создает и настраивает HTTP сервер
func initHTTPServer(config *Config, logg app.Logger, calendar *app.App) server.CalculatorServer {
	return internalhttp.NewServer(logg, config.Server.Host, config.Server.Port, calendar)
}

// createShutdownContext создает контекст для graceful shutdown
func createShutdownContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
}

// gracefulShutdown обрабатывает graceful shutdown HTTP сервера
func gracefulShutdown(ctx context.Context, server server.CalculatorServer, logg app.Logger) {
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}
}

// startGRPCServer запускает gRPC сервер
func startGRPCServer(config *Config, logg app.Logger, calendar *app.App, ctx context.Context) {
	grpcServer := internalgrpc.NewCalendarGRPCServer(logg, config.Server.GRPCPort, calendar)
	grpcServer.Start(ctx)
}
