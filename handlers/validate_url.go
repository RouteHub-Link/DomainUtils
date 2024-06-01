package handlers

import (
	"net/http"

	"github.com/RouteHub-Link/DomainUtils/tasks"
	"github.com/RouteHub-Link/DomainUtils/tasks/handler"
	"github.com/labstack/echo/v4"
)

type URLValidationHandlers struct {
	TaskServerConfig  *tasks.TaskServerConfig
	URLValidationTask *handler.URLValidationTask
}

func (dvh URLValidationHandlers) BindHandlers(e *echo.Echo) {
	e.POST("/validate/url", dvh.HandlePostValidateURL)
	e.GET("/validate/url/:id", dvh.HandleGetValidateURL)
}

func (dvh URLValidationHandlers) HandlePostValidateURL(c echo.Context) error {
	client := dvh.TaskServerConfig.NewClient()
	defer client.Close()

	validationPaylod := new(handler.URLValidationPayload)
	if err := c.Bind(validationPaylod); err != nil {
		return err
	}

	task, err := dvh.URLValidationTask.NewURLValidationTask(validationPaylod.Link)
	if err != nil {
		return err
	}

	info, err := client.Enqueue(task)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info.ID)
}

func (dvh URLValidationHandlers) HandleGetValidateURL(c echo.Context) error {
	inspector := dvh.TaskServerConfig.NewInspector()
	defer inspector.Close()

	infoPayload := new(handler.InfoPayload)
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "ID is required")
	}

	infoPayload.ID = id

	taskInfo, err := inspector.GetTaskInfo(dvh.URLValidationTask.Settings.Queue, infoPayload.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, taskInfo)
}
