package tasks

import (
	"log"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

type TaskServer struct {
	RedisAddr      string
	MonitoringPath string
	MonitoringPort string
}

func (t TaskServer) NewInspector() *asynq.Inspector {
	return asynq.NewInspector(asynq.RedisClientOpt{Addr: t.RedisAddr})
}

func (t TaskServer) NewClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: t.RedisAddr})
}

func (t TaskServer) Serve() {
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
		asynq.RedisClientOpt{Addr: t.RedisAddr},
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

func (t TaskServer) AsynqmonServe() {
	h := asynqmon.New(asynqmon.Options{
		RootPath:     t.MonitoringPath, // RootPath specifies the root for asynqmon app
		RedisConnOpt: asynq.RedisClientOpt{Addr: t.RedisAddr},
	})

	// Note: We need the tailing slash when using net/http.ServeMux.
	http.Handle(h.RootPath()+"/", h)
	println("Monitoring server is running link: http://localhost:" + t.MonitoringPort + t.MonitoringPath)

	// Go to http://localhost:8080/monitoring to see asynqmon homepage.
	log.Fatal(http.ListenAndServe(":"+t.MonitoringPort, nil))
}
