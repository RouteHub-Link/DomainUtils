package handlers

import (
	"net/http"

	"github.com/RouteHub-Link/DomainUtils/config"
	"github.com/RouteHub-Link/DomainUtils/tasks"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	ApplicationConfig *config.ApplicationConfig
	TaskServer        *tasks.TaskServer
}

func NewEchoServer(config *config.ApplicationConfig, taskServer *tasks.TaskServer) *EchoServer {
	return &EchoServer{
		ApplicationConfig: config,
		TaskServer:        taskServer,
	}
}

func (es EchoServer) Serve() {
	e := echo.New()
	config := es.ApplicationConfig

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if config.Health {
		e.GET("/health", func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})
	}

	url_validation_handlers := URLValidationHandlers{
		TaskServerConfig:  &config.TaskServerConfig,
		URLValidationTask: es.TaskServer.URLValidationTask,
	}

	dns_validation_handlers := DNSValidationHandlers{
		TaskServerConfig:  &config.TaskServerConfig,
		DNSValidationTask: es.TaskServer.DNSValidationTask,
	}

	site_validation_handlers := SiteValidationHandlers{
		TaskServerConfig:   &config.TaskServerConfig,
		SiteValidationTask: es.TaskServer.SiteValidationTask,
	}

	url_validation_handlers.BindHandlers(e)
	dns_validation_handlers.BindHandlers(e)
	site_validation_handlers.BindHandlers(e)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
