package config

type TaskConfigDNSValidation struct {
	TaskName     string `koanf:"task_name"`
	DNSTXTRecord string `koanf:"dns_txt_record"`
	TaskQueue    string `koanf:"task_queue"`
	TaskPriority int    `koanf:"task_priority"`
	DNSServer    string `koanf:"dns_server"`
}

var DefaultDNSValidationTaskConfig = TaskConfigDNSValidation{
	TaskName:     "dns:validate",
	DNSTXTRecord: "routehub_domainkey",
	TaskQueue:    "dns-validation",
	DNSServer:    "1.1.1.1:53",
	TaskPriority: 4,
}
