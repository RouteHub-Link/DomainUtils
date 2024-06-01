package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/RouteHub-Link/DomainUtils/tasks/config"
	"github.com/RouteHub-Link/DomainUtils/validator"

	"github.com/hibiken/asynq"
)

type URLValidationTask struct {
	Settings   config.Settings
	TaskConfig config.TaskConfigURLValidation
	Validator  *validator.Validator
}

type URLValidationPayload struct {
	Link string `json:"link"`
}

func (t *URLValidationTask) HandleURLValidationTask(ctx context.Context, task *asynq.Task) error {
	log.Printf("Processing task: %s, payload: %s", task.Type(), task.Payload())

	var payload URLValidationPayload
	var taskResultPayload *TaskResultPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		taskResultPayload = TaskResultPayload{}.New(false, "payload json unmarshal failed", err.Error())
		taskResultPayload.WriteResult(task)

		return nil
	}

	isValid, err := t.Validator.ValidateURL(payload.Link)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	taskResultPayload = TaskResultPayload{}.New(isValid, "", errMsg)
	taskResultPayload.WriteResult(task)

	return nil
}

func (t *URLValidationTask) NewURLValidationTask(link string) (*asynq.Task, error) {
	payload, err := json.Marshal(URLValidationPayload{
		Link: link,
	})

	return asynq.NewTask(t.TaskConfig.TaskName,
		payload,
		asynq.MaxRetry(t.Settings.MaxRetry),
		asynq.Timeout(t.Settings.Timeout),
		asynq.Queue(t.Settings.Queue),
		asynq.Deadline(time.Now().Add(t.Settings.DeadlineTimeout)),
		asynq.Retention(t.Settings.Retention),
	), err
}

func NewURLValidationTaskWithDefaults(URLValidationTaskConfig config.TaskConfigURLValidation, validator *validator.Validator) *URLValidationTask {
	return &URLValidationTask{
		Settings:   config.DefaultSettings(URLValidationTaskConfig.TaskQueue, URLValidationTaskConfig.TaskPriority),
		TaskConfig: URLValidationTaskConfig,
		Validator:  validator,
	}
}
