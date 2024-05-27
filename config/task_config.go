package config

import "github.com/RouteHub-Link/DomainUtils/tasks"

type TaskConfigs struct {
	DNSValidationTaskConfig tasks.DNSValidationTaskConfig `koanf:"dns_validation_task"`
	URLValidationTaskConfig tasks.URLValidationTaskConfig `koanf:"url_validation_task"`
}
