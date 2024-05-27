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

	domain_validation_handlers := DomainValidationHandlers{
		TaskServerConfig:  &config.TaskServerConfig,
		URLValidationTask: tasks.NewURLValidationTaskWithConfig(config.TaskConfigs.URLValidationTaskConfig),
	}

	dns_validation_handlers := DNSValidationHandlers{
		TaskServerConfig:  &config.TaskServerConfig,
		DNSValidationTask: tasks.NewDNSValidationTaskWithConfig(config.TaskConfigs.DNSValidationTaskConfig),
	}

	domain_validation_handlers.BindHandlers(e)
	dns_validation_handlers.BindHandlers(e)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
