package main

import (
	"net/http"

	"github.com/RouteHub-Link/DomainUtils/handlers"
	"github.com/RouteHub-Link/DomainUtils/tasks"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	_applicationConfig = GetApplicationConfig()
	URLValidationTask  = tasks.NewURLValidationTaskWithConfig(_applicationConfig.TaskConfigs.URLValidationTaskConfig)
	DNSValidationTask  = tasks.NewDNSValidationTaskWithConfig(_applicationConfig.TaskConfigs.DNSValidationTaskConfig)

	taskServer = &tasks.TaskServer{
		Config:            _applicationConfig.TaskServerConfig,
		DNSValidationTask: DNSValidationTask,
		URLValidationTask: URLValidationTask,
	}
)

func main() {
	go taskServer.Serve()
	go taskServer.AsynqmonServe()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if _applicationConfig.Health {
		e.GET("/health", func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})
	}

	domain_validation_handlers := handlers.DomainValidationHandlers{
		TaskServer: taskServer,
	}

	dns_validation_handlers := handlers.DNSValidationHandlers{
		TaskServer: taskServer,
	}

	domain_validation_handlers.BindHandlers(e)
	dns_validation_handlers.BindHandlers(e)

	e.Logger.Fatal(e.Start(":" + _applicationConfig.Port))
}
