package config

type TaskConfigURLValidation struct {
	TaskName     string `koanf:"task_name"`
	TaskQueue    string `koanf:"task_queue"`
	TaskPriority int    `koanf:"task_priority"`
}

var DefaultURLValidationTaskConfig = TaskConfigURLValidation{
	TaskName:     "url:validate",
	TaskQueue:    "url-validation",
	TaskPriority: 3,
}
