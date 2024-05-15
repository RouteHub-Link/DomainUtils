package handlers

import (
	"net/http"

	"github.com/RouteHub-Link/DomainUtils/tasks"
	"github.com/labstack/echo/v4"
)

type DomainValidationHandlers struct {
	TaskServer *tasks.TaskServer
}

func (dvh DomainValidationHandlers) BindHandlers(e *echo.Echo) {
	e.POST("/validate/url", dvh.HandlePostValidateURL)
	e.GET("/validate/url/:id", dvh.HandleGetValidateURL)
}

func (dvh DomainValidationHandlers) HandlePostValidateURL(c echo.Context) error {
	client := dvh.TaskServer.NewClient()
	defer client.Close()

	validationPaylod := new(tasks.URLValidationPayload)
	if err := c.Bind(validationPaylod); err != nil {
		return err
	}

	task, err := dvh.TaskServer.URLValidationTask.NewURLValidationTask(validationPaylod.Link)
	if err != nil {
		return err
	}

	info, err := client.Enqueue(task)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info.ID)
}

func (dvh DomainValidationHandlers) HandleGetValidateURL(c echo.Context) error {
	inspector := dvh.TaskServer.NewInspector()
	defer inspector.Close()

	infoPayload := new(tasks.InfoPayload)
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "ID is required")
	}

	infoPayload.ID = id

	taskInfo, err := inspector.GetTaskInfo(dvh.TaskServer.URLValidationTask.Settings.Queue, infoPayload.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, taskInfo)
}
