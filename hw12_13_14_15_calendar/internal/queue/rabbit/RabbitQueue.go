package rabbit

import (
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/queue"
)

type RabbitQueue struct {
}

func (q RabbitQueue) put(queue string, exchange string, message queue.MessageQueue) error {
	//TODO implement me
	panic("implement me")
}

func (q RabbitQueue) get(queue string) (queue.MessageQueue, error) {
	//TODO implement me
	panic("implement me")
}
