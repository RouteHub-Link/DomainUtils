package handlers

import (
	"net/http"

	"github.com/RouteHub-Link/DomainUtils/tasks"
	"github.com/labstack/echo/v4"
)

type DNSValidationHandlers struct {
	TaskServer        *tasks.TaskServer
	DNSValidationTask *tasks.DNSValidationTask
}

func (dvh DNSValidationHandlers) HandlePostValidateDNS(c echo.Context) error {
	client := dvh.TaskServer.NewClient()
	defer client.Close()

	validationPaylod := new(tasks.DNSValidationPayload)
	if err := c.Bind(validationPaylod); err != nil {
		return err
	}

	task, err := dvh.DNSValidationTask.NewURLValidationTask(validationPaylod.Link, validationPaylod.Value)
	if err != nil {
		return err
	}

	info, err := client.Enqueue(task)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info.ID)

}

func (dvh DNSValidationHandlers) HandleGetValidateDNS(c echo.Context) error {
	inspector := dvh.TaskServer.NewInspector()
	defer inspector.Close()

	infoPayload := new(tasks.InfoPayload)
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "ID is required")
	}

	infoPayload.ID = id

	taskInfo, err := inspector.GetTaskInfo(dvh.DNSValidationTask.Settings.Queue, infoPayload.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, taskInfo)
}

func (dvh DNSValidationHandlers) BindHandlers(e *echo.Echo) {
	e.POST("/validate/dns", dvh.HandlePostValidateDNS)
	e.GET("/validate/dns/:id", dvh.HandleGetValidateDNS)
}
