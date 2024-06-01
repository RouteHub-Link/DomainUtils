package tasks

import (
	"log"
	"net/http"

	"github.com/RouteHub-Link/DomainUtils/tasks/config"
	"github.com/RouteHub-Link/DomainUtils/tasks/handler"
	"github.com/RouteHub-Link/DomainUtils/validator"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

type TaskServer struct {
	Config             TaskServerConfig
	DNSValidationTask  *handler.DNSValidationTask
	URLValidationTask  *handler.URLValidationTask
	SiteValidationTask *handler.SiteValidationTask
}

func NewDefaultTaskServer(config TaskServerConfig, taskConfigs config.TaskConfigs, validator *validator.Validator) *TaskServer {
	return &TaskServer{
		Config:             config,
		DNSValidationTask:  handler.NewDNSValidationTaskWithDefaults(taskConfigs.DNSValidation, validator),
		URLValidationTask:  handler.NewURLValidationTaskWithDefaults(taskConfigs.URLValidation, validator),
		SiteValidationTask: handler.NewSiteValidationTaskWithDefaults(taskConfigs.SiteValidation, validator),
	}
}

func (t *TaskServer) Serve() {
	asynqQueue := map[string]int{
		"critical": 6,
		"default":  3,
		"low":      1,
	}

	asynqQueue[t.URLValidationTask.Settings.Queue] = t.URLValidationTask.Settings.QueuePriority
	asynqQueue[t.DNSValidationTask.Settings.Queue] = t.DNSValidationTask.Settings.QueuePriority
	asynqQueue[t.SiteValidationTask.Settings.Queue] = t.SiteValidationTask.Settings.QueuePriority

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
	mux.HandleFunc(t.SiteValidationTask.TaskConfig.TaskName, t.SiteValidationTask.HandleSiteValidationTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func (t *TaskServer) AsynqmonServe() {
	h := asynqmon.New(asynqmon.Options{
		RootPath:     t.Config.MonitoringPath,
		RedisConnOpt: asynq.RedisClientOpt{Addr: t.Config.RedisAddr},
	})

	http.Handle(h.RootPath()+"/", h)
	http.Handle("/", http.RedirectHandler(h.RootPath(), http.StatusFound))

	log.Printf("Monitoring server is running link: http://localhost:%s%s", t.Config.MonitoringPort, t.Config.MonitoringPath)
	log.Fatal(http.ListenAndServe(":"+t.Config.MonitoringPort, nil))
}
