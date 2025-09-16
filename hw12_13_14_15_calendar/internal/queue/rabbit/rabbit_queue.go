package rabbit

import (
	"context"
	"fmt"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitQueue[T any] struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	logger     app.Logger
	codec      Codec
}

func NewRabbitQueue[T any](url, username, password string, logger app.Logger) (queue.Queue[T], error) {
	connection, channel, err := connectToRabbitMQ(url, username, password)
	if err != nil {
		return nil, err
	}
	return &RabbitQueue[T]{
		connection: connection,
		channel:    channel,
		logger:     logger,
		codec:      JSONCodec{},
	}, nil
}

func connectToRabbitMQ(url, username, password string) (*amqp.Connection, *amqp.Channel, error) {
	// If url already contains the full AMQP URL, use it directly
	// Otherwise, construct it from components
	var amqpURL string
	if len(url) > 5 && url[:5] == "amqp:" {
		amqpURL = url
	} else {
		amqpURL = fmt.Sprintf("amqp://%s:%s@%s/", username, password, url)
	}
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}

func ensureQueue(ch *amqp.Channel, queueName, exchangeName, routingKey string) error {
	// Создаем эксчейндж если нужен
	if exchangeName != "" {
		err := ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
		if err != nil {
			return err
		}
	}

	// Создаем очередь
	_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Привязываем если есть эксчейндж
	if exchangeName != "" {
		err = ch.QueueBind(queueName, routingKey, exchangeName, false, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (q RabbitQueue[T]) Put(queue string, exchange string, message queue.MessageQueue[T]) error {
	routingKey := "default"
	if err := ensureQueue(q.channel, queue, exchange, routingKey); err != nil {
		return err
	}
	body, err := q.codec.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return q.channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: q.codec.ContentType(),
		MessageId:   message.ID,
		Body:        body,
	})
}

func (q RabbitQueue[T]) Get(context context.Context, queueName string) (<-chan queue.MessageQueue[T], <-chan error) {
	ch := make(chan queue.MessageQueue[T])
	errch := make(chan error)
	go func() {
		q.get(context, ch, errch, queueName)
	}()
	return ch, errch
}

func (q RabbitQueue[T]) get(ctx context.Context, resCh chan queue.MessageQueue[T], errCh chan error, queueName string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			d, ok, err := q.channel.Get(queueName, false)
			if err != nil {
				errCh <- fmt.Errorf("get message: %w", err)
				continue
			}
			if !ok {
				//q.logger.Info("No messages")
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Декодируем сообщение
			var message queue.MessageQueue[T]
			if err := q.codec.Unmarshal(d.Body, &message); err != nil {
				errCh <- fmt.Errorf("unmarshal: %w", err)
				d.Ack(false)
				continue
			}

			// Устанавливаем ID из сообщения RabbitMQ
			if message.ID == "" {
				message.ID = d.MessageId
			}
			resCh <- message
			d.Ack(false)
		}
	}
}

func (q RabbitQueue[T]) Close() error {
	q.channel.Close()
	return q.connection.Close()
}
