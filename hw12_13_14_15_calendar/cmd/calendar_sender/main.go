package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
)

func main() {
	configPath := flag.String("config", "./configs/sender_config.yaml", "Path to config file")
	flag.Parse()

	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logg := logger.New(config.Logger.Level)

	queue, err := rabbit.NewRabbitQueue[storage.Notification](
		config.Rabbit.URL,
		config.Rabbit.Username,
		config.Rabbit.Password,
		logg,
	)
	if err != nil {
		logg.Error("Failed to connect to RabbitMQ: " + err.Error())
		os.Exit(1)
	}
	defer queue.Close()

	logg.Info("Sender started")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msgChan, errChan := queue.Get(ctx, config.EventQueue.Name)

				for {
					select {
					case <-ctx.Done():
						return
					case msg := <-msgChan:
						notification := msg.Body
						logg.Info(fmt.Sprintf("Sending notification: EventID=%s, Title=%s, UserID=%s, EventTime=%s",
							notification.EventID, notification.Title, notification.UserID, notification.EventTime.Format("2006-01-02 15:04:05")))
					case err := <-errChan:
						if err != nil {
							logg.Error("Error receiving message: " + err.Error())
							return
						}
					}
				}
			}
		}
	}()

	<-sigChan
	logg.Info("Sender stopped")
}
