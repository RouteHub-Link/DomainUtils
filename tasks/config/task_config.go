package config

type TaskConfigs struct {
	DNSValidation  TaskConfigDNSValidation  `koanf:"dns_validation_task"`
	URLValidation  TaskConfigURLValidation  `koanf:"url_validation_task"`
	SiteValidation TaskConfigSiteValidation `koanf:"site_validation_task"`
}
