package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/RouteHub-Link/DomainUtils/validator"
	"github.com/hibiken/asynq"
)

const (
	TaskTypeURLValidate              = "url:validate"
	TaskTypeURLValidateQueue         = "url-validation"
	TaskTypeURLValidateQueuePriority = 3
)

type URLValidationTask struct {
	Settings Settings
}

type URLValidationPayload struct {
	Link string `json:"link"`
}

func NewURLValidationTask(settings Settings) *URLValidationTask {
	return &URLValidationTask{
		Settings: settings,
	}
}

func NewURLValidationTaskWithDefaults() *URLValidationTask {
	return &URLValidationTask{
		Settings: DefaultSettings(TaskTypeURLValidateQueue, TaskTypeURLValidateQueuePriority),
	}
}

func (t *URLValidationTask) NewURLValidationTask(link string) (*asynq.Task, error) {
	payload, err := json.Marshal(URLValidationPayload{
		Link: link,
	})

	return asynq.NewTask(TaskTypeURLValidate,
		payload,
		asynq.MaxRetry(t.Settings.MaxRetry),
		asynq.Timeout(t.Settings.Timeout),
		asynq.Queue(t.Settings.Queue),
		asynq.Deadline(time.Now().Add(t.Settings.DeadlineTimeout)),
		asynq.Retention(t.Settings.Retention),
	), err
}

func (t *URLValidationTask) HandleURLValidationTask(ctx context.Context, task *asynq.Task) error {
	log.Printf("Processing task: %s, payload: %s", task.Type(), task.Payload())

	var payload URLValidationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("payload json unmarshal failed: %v: %w", err, asynq.SkipRetry)

	}

	isValid, err := _validator.ValidateURL(payload.Link)
	URLTaskResultPayload := TaskResultPayload{
		IsValid: isValid,
	}

	log.Printf("\nValidation result: %v \n Err : %v", URLTaskResultPayload, err)

	if err != nil {
		if ce, ok := err.(*validator.CheckError); ok {
			errMsg := err.Error()

			// Try agein if the error is unreachable
			if ce.Msg == validator.CheckErrorMessages[validator.ErrUnreachable] {
				return err
			}

			URLTaskResultPayload.Error = &errMsg
			payloadJson, _ := URLTaskResultPayload.ToJson()

			_, _ = task.ResultWriter().Write(payloadJson)

			return fmt.Errorf("validation failed: %v: %w", err, asynq.SkipRetry)
		} else {
			return err
		}
	}

	payloadJson, _ := URLTaskResultPayload.ToJson()
	_, _ = task.ResultWriter().Write(payloadJson)

	return err
}
