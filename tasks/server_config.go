package tasks

import "github.com/hibiken/asynq"

type TaskServerConfig struct {
	RedisAddr      string `koanf:"redis_addr"`
	MonitoringDash bool   `koanf:"monitoring_dashboard"`
	MonitoringPath string `koanf:"monitoring_path"`
	MonitoringPort string `koanf:"monitoring_port"`
	Concurrency    int    `koanf:"concurrency"`
}

func (t *TaskServerConfig) NewInspector() *asynq.Inspector {
	return asynq.NewInspector(asynq.RedisClientOpt{Addr: t.RedisAddr})
}

func (t *TaskServerConfig) NewClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: t.RedisAddr})
}

var DefaultTaskServerConfig = TaskServerConfig{
	RedisAddr:      "localhost:6379",
	MonitoringDash: true,
	MonitoringPath: "/monitoring",
	MonitoringPort: "8081",
	Concurrency:    10,
}
