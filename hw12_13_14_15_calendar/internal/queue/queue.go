package queue

import (
	"context"
)

type MessageQueue[T any] struct {
	ID   string
	Body T
}

type Queue[T any] interface {
	Put(queue string, exchange string, message MessageQueue[T]) error
	Get(context context.Context, queueName string) (<-chan MessageQueue[T], <-chan error)
	Close() error
}
