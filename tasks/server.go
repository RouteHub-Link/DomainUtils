package tasks

import (
	"log"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

type TaskServer struct {
	Config            TaskServerConfig
	DNSValidationTask *DNSValidationTask
	URLValidationTask *URLValidationTask
}

type TaskServerConfig struct {
	RedisAddr      string `koanf:"redis_addr"`
	MonitoringPath string `koanf:"monitoring_path"`
	MonitoringPort string `koanf:"monitoring_port"`
	Concurrency    int    `koanf:"concurrency"`
}

var DefaultTaskServerConfig = TaskServerConfig{
	RedisAddr:      "localhost:6379",
	MonitoringPath: "/monitoring",
	MonitoringPort: "8081",
	Concurrency:    10,
}

func (t *TaskServer) NewInspector() *asynq.Inspector {
	return asynq.NewInspector(asynq.RedisClientOpt{Addr: t.Config.RedisAddr})
}

func (t *TaskServer) NewClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: t.Config.RedisAddr})
}

func (t *TaskServer) Serve() {
	asynqQueue := map[string]int{
		"critical": 6,
		"default":  3,
		"low":      1,
	}

	asynqQueue[t.URLValidationTask.Settings.Queue] = t.URLValidationTask.Settings.QueuePriority
	asynqQueue[t.DNSValidationTask.Settings.Queue] = t.DNSValidationTask.Settings.QueuePriority

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: t.Config.RedisAddr},
		asynq.Config{
			Concurrency: t.Config.Concurrency,
			Queues:      asynqQueue,
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(t.URLValidationTask.TaskConfig.TaskName, t.URLValidationTask.HandleURLValidationTask)
	mux.HandleFunc(t.DNSValidationTask.TaskConfig.TaskName, t.DNSValidationTask.HandleDNSValidationTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func (t *TaskServer) AsynqmonServe() {
	h := asynqmon.New(asynqmon.Options{
		RootPath:     t.Config.MonitoringPath, // RootPath specifies the root for asynqmon app
		RedisConnOpt: asynq.RedisClientOpt{Addr: t.Config.RedisAddr},
	})

	// Note: We need the tailing slash when using net/http.ServeMux.
	http.Handle(h.RootPath()+"/", h)
	log.Printf("Monitoring server is running link: http://localhost:%s%s", t.Config.MonitoringPort, t.Config.MonitoringPath)

	// Go to http://localhost:8080/monitoring to see asynqmon homepage.
	log.Fatal(http.ListenAndServe(":"+t.Config.MonitoringPort, nil))
}
