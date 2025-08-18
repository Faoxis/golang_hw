package queue

type MessageQueue struct {
	ID   string
	body any
}

type Queue interface {
	put(queue string, exchange string, message MessageQueue) error
	get(queue string) (MessageQueue, error)
}
