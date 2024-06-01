package config

type TaskConfigSiteValidation struct {
	TaskName     string `koanf:"task_name"`
	TaskQueue    string `koanf:"task_queue"`
	TaskPriority int    `koanf:"task_priority"`
}

var DefaultSiteValidationTaskConfig = TaskConfigSiteValidation{
	TaskName:     "site:validate",
	TaskQueue:    "site-validation",
	TaskPriority: 3,
}
