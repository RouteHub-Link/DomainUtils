package handlers

import (
	"net/http"

	"github.com/RouteHub-Link/DomainUtils/tasks"
	"github.com/RouteHub-Link/DomainUtils/tasks/handler"
	"github.com/labstack/echo/v4"
)

type SiteValidationHandlers struct {
	TaskServerConfig   *tasks.TaskServerConfig
	SiteValidationTask *handler.SiteValidationTask
}

func (svh SiteValidationHandlers) BindHandlers(e *echo.Echo) {
	e.POST("/validate/site", svh.HandlePostValidateSite)
	e.GET("/validate/site/:id", svh.HandleGetValidateSite)
}

func (svh SiteValidationHandlers) HandlePostValidateSite(c echo.Context) error {
	client := svh.TaskServerConfig.NewClient()
	defer client.Close()

	validationPaylod := new(handler.SiteValidationPayload)
	if err := c.Bind(validationPaylod); err != nil {
		return err
	}

	task, err := svh.SiteValidationTask.NewSiteValidationTask(validationPaylod.Link)
	if err != nil {
		return err
	}

	info, err := client.Enqueue(task)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info.ID)
}

func (svh SiteValidationHandlers) HandleGetValidateSite(c echo.Context) error {
	inspector := svh.TaskServerConfig.NewInspector()
	defer inspector.Close()

	infoPayload := new(handler.InfoPayload)
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "ID is required")
	}

	infoPayload.ID = id

	taskInfo, err := inspector.GetTaskInfo(svh.SiteValidationTask.Settings.Queue, infoPayload.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, taskInfo)
}
