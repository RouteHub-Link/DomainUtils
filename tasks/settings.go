package tasks

import (
	"encoding/json"
	"time"

	"github.com/RouteHub-Link/URLValidator/validator"
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

type InfoPayload struct {
	ID string `json:"id"`
}

type TaskResultPayload struct {
	IsValid bool    `json:"isValid"`
	Error   *string `json:"error,omitempty"`
}

func (p TaskResultPayload) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p TaskResultPayload) FromJson(data []byte) error {
	return json.Unmarshal(data, &p)
}

var _validator = validator.DefaultValidator()
