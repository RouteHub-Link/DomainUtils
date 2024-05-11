package tasks

import (
	"log"

	"github.com/hibiken/asynq"
)

func Serve(redisAddr string) {
	URLValidationTask := NewURLValidationTaskWithDefaults()
	DNSValidationTask := NewDNSValidationTaskWithDefaults()

	asynqQueue := map[string]int{
		"critical": 6,
		"default":  3,
		"low":      1,
	}

	asynqQueue[URLValidationTask.Settings.Queue] = URLValidationTask.Settings.QueuePriority
	asynqQueue[DNSValidationTask.Settings.Queue] = DNSValidationTask.Settings.QueuePriority

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues:      asynqQueue,
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskTypeURLValidate, URLValidationTask.HandleURLValidationTask)
	mux.HandleFunc(TaskTypeDNSValidate, DNSValidationTask.HandleDNSValidationTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
