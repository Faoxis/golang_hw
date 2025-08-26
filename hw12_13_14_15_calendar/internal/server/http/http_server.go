package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/server"

	"github.com/go-chi/chi/v5/middleware"

	route "github.com/go-chi/chi/v5"
)

type HttpServer struct {
	logger      server.Logger
	application server.Application
	server      *http.Server
}

func NewServer(logger server.Logger, host string, port int, app server.Application) server.CalculatorServer {
	router := route.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(func(next http.Handler) http.Handler {
		return loggingMiddleware(logger, next)
	})

	router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
	})
	router.Route("/events", func(router route.Router) {
		router.Get("/day", getEventsForDay(app, logger))
		router.Get("/week", getEventsForWeek(app, logger))
		router.Get("/month", getEventsForMonth(app, logger))
		router.Get("/", getEvents(app, logger)) // старый универсальный, можно оставить для обратной совместимости
		router.Post("/", addNewEvent(app, logger))
		router.Get("/{id}", getEvent(app, logger))
		router.Put("/{id}", updateEvent(app, logger))
		router.Delete("/{id}", deleteEvent(app, logger))
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	return &HttpServer{
		logger:      logger,
		application: app,
		server:      srv,
	}
}

func (s *HttpServer) Start(ctx context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("http server error: " + err.Error())
		}
	}()
	<-ctx.Done()
	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)

}
