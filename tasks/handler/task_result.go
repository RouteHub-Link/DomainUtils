package handler

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

type InfoPayload struct {
	ID string `json:"id"`
}

type TaskResultPayload struct {
	IsValid bool    `json:"isValid"`
	Message *string `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}

func (p TaskResultPayload) New(isValid bool, message string, err string) *TaskResultPayload {
	_message := &message
	_err := &err

	return &TaskResultPayload{
		IsValid: isValid,
		Message: _message,
		Error:   _err,
	}
}

func (p TaskResultPayload) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p TaskResultPayload) FromJson(data []byte) error {
	return json.Unmarshal(data, &p)
}

func (p TaskResultPayload) WriteResult(task *asynq.Task) {
	payloadJson, _ := p.ToJson()
	task.ResultWriter().Write(payloadJson)
}
