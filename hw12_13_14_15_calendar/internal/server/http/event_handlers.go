package internalhttp

import (
	"fmt"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func getEvents(application Application, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		var events []storage.Event
		var err error
		if day := query.Get("day"); day != "" {
			date, err := time.Parse("2006-01-02", day)
			if err != nil {
				logger.Warn(fmt.Sprintf("invalid date format: %s", day))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			events, err = application.ListEventsForDay(r.Context(), date)
		} else if week := query.Get("week"); week != "" {
			date, err := time.Parse("2006-01-02", week)
			if err != nil {
				logger.Warn(fmt.Sprintf("invalid date format: %s, err %s", week, err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			events, err = application.ListEventsForWeek(r.Context(), date)
		} else if month := query.Get("month"); month != "" {
			date, err := time.Parse("2006-01-02", month)
			if err != nil {
				logger.Warn(fmt.Sprintf("invalid date format: %s", month))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			events, err = application.ListEventsForMonth(r.Context(), date)
		} else {
			logger.Warn(fmt.Sprintf("invalid query: %s", query))
			w.WriteHeader(http.StatusBadRequest)
		}
		if err != nil {
			logger.Error(fmt.Sprintf("error getting events: %s", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseEvents := make([]EventResponse, 0, len(events))
		for _, event := range events {
			responseEvents = append(responseEvents, mapStorageEventToEventResponse(event))
		}
		err = sendInResponse(w, responseEvents, http.StatusOK)
		if err != nil {
			logger.Warn(fmt.Sprintf("send response error: %s", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func deleteEvent(application Application, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		err := application.DeleteEvent(r.Context(), id)
		if err != nil {
			logger.Warn(fmt.Errorf("can't delete event: %w", err).Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func updateEvent(application Application, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		eventRequest, err := fromJson[EventRequest](r.Body)
		if err != nil {
			logger.Warn(fmt.Sprintf("error parsing event: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = application.UpdateEvent(
			r.Context(),
			id, eventRequest.Title, eventRequest.Description, eventRequest.UserID,
			eventRequest.StartTime,
			eventRequest.Duration, eventRequest.NotifyBefore,
		)
		updatedStorageEvent, err := application.GetEventByID(r.Context(), id)
		if err != nil {
			logger.Warn(fmt.Sprintf("error updating event: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = sendInResponse(w, updatedStorageEvent, http.StatusOK)
		if err != nil {
			logger.Warn(fmt.Sprintf("error writing response: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func getEvent(app Application, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		event, err := app.GetEventByID(r.Context(), id)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to get event id: %v", err))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = sendInResponse(
			w,
			event,
			http.StatusOK,
		)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to write response: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func addNewEvent(app Application, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := fromJson[EventRequest](r.Body)
		if err != nil {
			logger.Warn(fmt.Sprintf("error parsing event: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newUUID, err := uuid.NewUUID()
		if err != nil {
			logger.Warn("Request ID cannot be parsed")
			w.WriteHeader(http.StatusBadRequest)
		}
		id := newUUID.String()

		err = app.CreateEvent(
			r.Context(),
			id, event.Title, event.Description, event.UserID,
			event.StartTime, event.Duration, event.NotifyBefore,
		)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to create event: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		savedEvent, err := app.GetEventByID(r.Context(), id)
		err = sendInResponse(
			w,
			mapStorageEventToEventResponse(savedEvent),
			http.StatusOK,
		)
		if err != nil {
			logger.Warn("Failed to write response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
