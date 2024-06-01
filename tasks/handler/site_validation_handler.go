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

type SiteValidationTask struct {
	Settings   config.Settings
	TaskConfig config.TaskConfigSiteValidation
	Validator  *validator.Validator
}

type SiteValidationPayload struct {
	Link string `json:"link"`
}

func (t *SiteValidationTask) HandleSiteValidationTask(ctx context.Context, task *asynq.Task) error {
	log.Printf("Processing task: %s, payload: %s", task.Type(), task.Payload())

	var payload SiteValidationPayload
	var taskResultPayload *TaskResultPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		taskResultPayload = TaskResultPayload{}.New(false, "payload json unmarshal failed", err.Error())
		taskResultPayload.WriteResult(task)

		return nil
	}
	errMsg := ""

	isValid, err := t.Validator.ValidateSite(payload.Link)
	if err != nil {
		if err.Error() == validator.CheckErrorMessages[validator.ErrUnreachable] {
			return err
		}

		errMsg = err.Error()
	}

	taskResultPayload = TaskResultPayload{}.New(isValid, "", errMsg)
	taskResultPayload.WriteResult(task)

	return nil
}

func (t *SiteValidationTask) NewSiteValidationTask(link string) (*asynq.Task, error) {
	payload, err := json.Marshal(SiteValidationPayload{
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

func NewSiteValidationTaskWithDefaults(SiteValidationTaskConfig config.TaskConfigSiteValidation, validator *validator.Validator) *SiteValidationTask {
	return &SiteValidationTask{
		Settings:   config.DefaultSettings(SiteValidationTaskConfig.TaskQueue, SiteValidationTaskConfig.TaskPriority),
		TaskConfig: SiteValidationTaskConfig,
		Validator:  validator,
	}
}
