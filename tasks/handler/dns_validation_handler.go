package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/RouteHub-Link/DomainUtils/tasks/config"
	"github.com/RouteHub-Link/DomainUtils/validator"
	"github.com/hibiken/asynq"
)

type DNSValidationTask struct {
	Settings   config.Settings
	TaskConfig config.TaskConfigDNSValidation
	Validator  *validator.Validator
}

type DNSValidationPayload struct {
	Link  string `json:"link"`
	Value string `json:"value"`
}

func (t *DNSValidationTask) HandleDNSValidationTask(ctx context.Context, task *asynq.Task) error {
	log.Printf("Processing task: %s, payload: %s", task.Type(), task.Payload())
	var payload DNSValidationPayload

	var taskResultPayload *TaskResultPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		taskResultPayload = TaskResultPayload{}.New(false, "payload json unmarshal failed", err.Error())
		taskResultPayload.WriteResult(task)

		return nil
	}

	isvalid, err := t.Validator.ValidateOwnershipOverDNSTxtRecord(payload.Link, t.TaskConfig.DNSTXTRecord, payload.Value, t.TaskConfig.DNSServer)
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	taskResultPayload = TaskResultPayload{}.New(isvalid, "", errorMsg)

	// if the error is due to the domain being unreachable, return the error for retry
	if !isvalid {
		result := taskResultPayload.New(isvalid, "error validating payload", errorMsg)
		payloadBytes, _ := result.ToJson()
		return fmt.Errorf("%v", string(payloadBytes))
	}

	taskResultPayload.WriteResult(task)
	return nil
}

func (t *DNSValidationTask) NewDNSValidationTask(link string, value string) (*asynq.Task, error) {
	payload, err := json.Marshal(DNSValidationPayload{
		Link:  link,
		Value: value,
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

func NewDNSValidationTaskWithDefaults(DNSValidationTaskConfig config.TaskConfigDNSValidation, validator *validator.Validator) *DNSValidationTask {
	return &DNSValidationTask{
		Settings:   config.DefaultSettings(DNSValidationTaskConfig.TaskQueue, DNSValidationTaskConfig.TaskPriority),
		TaskConfig: DNSValidationTaskConfig,
		Validator:  validator,
	}
}
