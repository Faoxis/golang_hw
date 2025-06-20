package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	logger      Logger
	application Application
	server      *http.Server
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Warn(msg string)
}

type Application interface {
	CreateEvent(
		ctx context.Context,
		id, title, description, userID string,
		startTime time.Time,
		duration, notifyBefore calendar_types.CalendarDuration,
	) error

	UpdateEvent(
		ctx context.Context,
		id, title, description, userID string,
		startTime time.Time,
		duration, notifyBefore calendar_types.CalendarDuration,
	) error

	DeleteEvent(ctx context.Context, id string) error
	GetEventByID(ctx context.Context, id string) (storage.Event, error)
	ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListEventsForWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListEventsForMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func NewServer(logger Logger, host string, port int, app Application) *Server {
	router := chi.NewRouter()

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
	router.Route("/events", func(router chi.Router) {
		router.Get("/", getEvents(app, logger))
		router.Post("/", addNewEvent(app, logger))
		router.Get("/{id}", getEvent(app, logger))
		router.Put("/{id}", updateEvent(app, logger))
		router.Delete("/{id}", deleteEvent(app, logger))
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	return &Server{
		logger:      logger,
		application: app,
		server:      srv,
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("http server error: " + err.Error())
		}
	}()
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)

}
