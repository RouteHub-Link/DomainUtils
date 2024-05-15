package main

import (
	"net/http"

	"github.com/RouteHub-Link/DomainUtils/handlers"
	"github.com/RouteHub-Link/DomainUtils/tasks"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	taskServer        = tasks.TaskServer{RedisAddr: "127.0.0.1:6379", MonitoringPath: "/monitoring", MonitoringPort: "8080"}
	URLValidationTask = tasks.NewURLValidationTaskWithDefaults()
	DNSValidationTask = tasks.NewDNSValidationTaskWithDefaults()
)

func main() {
	go taskServer.Serve()
	go taskServer.AsynqmonServe()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	domain_validation_handlers := handlers.DomainValidationHandlers{
		TaskServer:        &taskServer,
		URLValidationTask: URLValidationTask,
	}

	dns_validation_handlers := handlers.DNSValidationHandlers{
		TaskServer:        &taskServer,
		DNSValidationTask: DNSValidationTask,
	}

	domain_validation_handlers.BindHandlers(e)
	dns_validation_handlers.BindHandlers(e)

	e.Logger.Fatal(e.Start(":1325"))
}
