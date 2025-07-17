package internalgrpc

import (
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/api"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func convertToEventResponse(storageEvent storage.Event) api.EventResponse {
	return api.EventResponse{
		Id:           storageEvent.ID,
		Title:        storageEvent.Title,
		StartTime:    timestamppb.New(storageEvent.StartTime),
		Description:  storageEvent.Description,
		Duration:     time.Duration(storageEvent.Duration).String(),
		UserId:       storageEvent.UserID,
		NotifyBefore: time.Duration(storageEvent.NotifyBefore).String(),
	}
}
