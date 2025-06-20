//go:build integration

package integration

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/"
)

func TestHelloEndpoint(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := internalhttp.NewServer(logger, "localhost", 8080, app)
	go func() {
		_ = srv.Start(ctx)
	}()

	// ждем, чтобы сервер успел подняться
	time.Sleep(500 * time.Millisecond)

	ctx := context.Background()

	// Запускаем PostgreSQL-контейнер
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_USER":     "calendar",
			"POSTGRES_DB":       "calendar",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}
	defer pgContainer.Terminate(ctx)

	// Получаем хост и порт
	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")
	t.Logf("PostgreSQL started on %s:%s", host, port.Port())

	// Тест предполагает, что сервер уже поднят и слушает localhost:8080
	time.Sleep(2 * time.Second) // Подождать запуск сервера, если нужно

	resp, err := http.Get("http://localhost:8080/hello")
	if err != nil {
		t.Fatalf("HTTP GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello, World!" {
		t.Errorf("unexpected body: %s", string(body))
	}
}
