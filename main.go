package main

import (
	"github.com/RouteHub-Link/DomainUtils/config"
	"github.com/RouteHub-Link/DomainUtils/handlers"
	"github.com/RouteHub-Link/DomainUtils/tasks"
)

var (
	applicationConfig = config.GetApplicationConfig()
	URLValidationTask = tasks.NewURLValidationTaskWithConfig(applicationConfig.TaskConfigs.URLValidationTaskConfig)
	DNSValidationTask = tasks.NewDNSValidationTaskWithConfig(applicationConfig.TaskConfigs.DNSValidationTaskConfig)

	taskServer = &tasks.TaskServer{
		Config:            applicationConfig.TaskServerConfig,
		DNSValidationTask: DNSValidationTask,
		URLValidationTask: URLValidationTask,
	}
)

func main() {
	switch applicationConfig.HostingMode {
	case config.TaskClient:
		es := handlers.EchoServer{
			ApplicationConfig: applicationConfig,
		}
		es.Serve()
	case config.TaskServer:
		taskServer.Serve()
	case config.TaskMonitoring:
		taskServer.AsynqmonServe()
	}
}
