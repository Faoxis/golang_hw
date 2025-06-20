package storage

import "time"

type Notification struct {
	EventID   string    // ID события
	Title     string    // Название события
	EventTime time.Time // Дата события
	UserID    string    // Кому отправить
}
