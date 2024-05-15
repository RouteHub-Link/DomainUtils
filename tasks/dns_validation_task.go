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
	TaskTypeDNSValidate          = "dns:validate"
	TaskTypeDNSValidateTXTRecord = "routehub_domainkey"
	TaskTypeDNSValidateQueue     = "dns-validation"
	TaskTypeDNSValidatePriority  = 4
)

type DNSValidationTask struct {
	Settings      Settings
	TXTRecordName string
}

type DNSValidationPayload struct {
	Link  string `json:"link"`
	Value string `json:"value"`
}

func NewDNSValidationTask(settings Settings) *DNSValidationTask {
	return &DNSValidationTask{
		Settings:      settings,
		TXTRecordName: TaskTypeDNSValidateTXTRecord,
	}
}

func NewDNSValidationTaskWithDefaults() *DNSValidationTask {
	return &DNSValidationTask{
		Settings:      DefaultSettings(TaskTypeDNSValidateQueue, TaskTypeDNSValidatePriority),
		TXTRecordName: TaskTypeDNSValidateTXTRecord,
	}
}

func (t *DNSValidationTask) NewURLValidationTask(link string, value string) (*asynq.Task, error) {
	payload, err := json.Marshal(DNSValidationPayload{
		Link:  link,
		Value: value,
	})

	return asynq.NewTask(TaskTypeDNSValidate,
		payload,
		asynq.MaxRetry(t.Settings.MaxRetry),
		asynq.Timeout(t.Settings.Timeout),
		asynq.Queue(t.Settings.Queue),
		asynq.Deadline(time.Now().Add(t.Settings.DeadlineTimeout)),
		asynq.Retention(t.Settings.Retention),
	), err
}

func (t *DNSValidationTask) HandleDNSValidationTask(ctx context.Context, task *asynq.Task) error {
	log.Printf("Processing task: %s, payload: %s", task.Type(), task.Payload())
	var payload DNSValidationPayload
	res := TaskResultPayload{}

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		result := TaskResultPayload{}.New(false, "payload json unmarshal failed", err.Error())
		payloadBytes, _ := result.ToJson()
		return fmt.Errorf("%v : %w", string(payloadBytes), asynq.SkipRetry)
	}

	isvalid, err := _validator.ValidateOwnershipOverDNSTxtRecord(payload.Link, t.TXTRecordName, payload.Value)
	res.IsValid = isvalid

	if err != nil {
		errorStr := err.Error()
		res.Error = &errorStr

		if err.Error() == validator.CheckErrorMessages[validator.ErrDNSNameValueNull] {

			result := res.New(isvalid, "error validating payload", errorStr)
			payloadBytes, _ := result.ToJson()

			return fmt.Errorf("%v: %w", string(payloadBytes), asynq.SkipRetry)
		} else {
			result := res.New(isvalid, "error validating ownership over DNS TXT record", errorStr)
			payloadBytes, _ := result.ToJson()

			return fmt.Errorf("%v", string(payloadBytes))
		}
	}

	payloadBytes, _ := res.ToJson()
	_, _ = task.ResultWriter().Write(payloadBytes)
	return nil
}
