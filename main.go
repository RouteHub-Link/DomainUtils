package main

import (
	"os"

	"net/http"

	"github.com/RouteHub-Link/URLValidator/tasks"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	redisAddr         = "127.0.0.1:6379"
	URLValidationTask = tasks.NewURLValidationTaskWithDefaults()
	DNSValidationTask = tasks.NewDNSValidationTaskWithDefaults()
)

func main() {
	redisAddr = os.Getenv("redisAddr")
	go tasks.Serve(redisAddr)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/validate/url/:id", HandleGetValidateURL)
	e.POST("/validate/url", HandlePostValidateURL)

	e.GET("/validate/dns/:id", HandleGetValidateDNS)
	e.POST("/validate/dns", HandlePostValidateDNS)

	e.Logger.Fatal(e.Start(":1325"))
}

func HandlePostValidateURL(c echo.Context) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	validationPaylod := new(tasks.URLValidationPayload)
	if err := c.Bind(validationPaylod); err != nil {
		return err
	}

	task, err := URLValidationTask.NewURLValidationTask(validationPaylod.Link)
	if err != nil {
		return err
	}

	info, err := client.Enqueue(task)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info.ID)
}

func HandleGetValidateURL(c echo.Context) error {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})
	defer inspector.Close()

	infoPayload := new(tasks.InfoPayload)
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "ID is required")
	}

	infoPayload.ID = id

	taskInfo, err := inspector.GetTaskInfo(URLValidationTask.Settings.Queue, infoPayload.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, taskInfo)
}

func HandlePostValidateDNS(c echo.Context) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	validationPaylod := new(tasks.DNSValidationPayload)
	if err := c.Bind(validationPaylod); err != nil {
		return err
	}

	task, err := DNSValidationTask.NewURLValidationTask(validationPaylod.Link, validationPaylod.Value)
	if err != nil {
		return err
	}

	info, err := client.Enqueue(task)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info.ID)

}

func HandleGetValidateDNS(c echo.Context) error {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})
	defer inspector.Close()

	infoPayload := new(tasks.InfoPayload)
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "ID is required")
	}

	infoPayload.ID = id

	taskInfo, err := inspector.GetTaskInfo(DNSValidationTask.Settings.Queue, infoPayload.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, taskInfo)
}
