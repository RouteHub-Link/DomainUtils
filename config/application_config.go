package config

import (
	"github.com/RouteHub-Link/DomainUtils/tasks"
	task_config "github.com/RouteHub-Link/DomainUtils/tasks/config"
	"github.com/RouteHub-Link/DomainUtils/validator"
)

type ApplicationConfig struct {
	ValidatorConfig  validator.CheckConfig   `koanf:"validator"`
	TaskServerConfig tasks.TaskServerConfig  `koanf:"task_server"`
	TaskConfigs      task_config.TaskConfigs `koanf:"tasks"`
	Port             string                  `koanf:"port"`
	Health           bool                    `koanf:"health"`
	HostingMode      HostingMode             `koanf:"hosting_mode"`
}
