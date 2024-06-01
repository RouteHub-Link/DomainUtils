package config

import (
	"time"
)

type Settings struct {
	MaxRetry        int
	Timeout         time.Duration
	DeadlineTimeout time.Duration
	Queue           string
	QueuePriority   int
	Retention       time.Duration
}

func (s Settings) GetPriority() map[string]int {
	return map[string]int{
		s.Queue: s.QueuePriority,
	}
}

func DefaultSettings(Queue string, QueuePriority int) Settings {
	return Settings{
		MaxRetry:        10,
		Timeout:         2 * time.Minute,
		DeadlineTimeout: 24 * time.Hour,
		Queue:           Queue,
		QueuePriority:   QueuePriority,
		Retention:       10 * (24 * time.Hour),
	}
}
