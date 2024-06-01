package main

import (
	"github.com/RouteHub-Link/DomainUtils/config"
	"github.com/RouteHub-Link/DomainUtils/handlers"
	"github.com/RouteHub-Link/DomainUtils/tasks"
	"github.com/RouteHub-Link/DomainUtils/validator"
)

var (
	applicationConfig = config.GetApplicationConfig()
	taskValidator     = validator.DefaultValidator()

	taskServer = tasks.NewDefaultTaskServer(
		applicationConfig.TaskServerConfig,
		applicationConfig.TaskConfigs,
		taskValidator,
	)
)

func main() {
	switch applicationConfig.HostingMode {
	case config.TaskReceiver:
		es := handlers.NewEchoServer(applicationConfig, taskServer)
		es.Serve()
	case config.TaskServer:
		taskServer.Serve()
	case config.TaskMonitoring:
		taskServer.AsynqmonServe()
	}
}
