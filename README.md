# DomainUtils

![gov120](https://github.com/RouteHub-Link/DomainUtils/assets/16222645/3d751aaa-ad97-45a4-9c89-99521729e44a)
[![Go Report Card](https://goreportcard.com/badge/github.com/RouteHub-Link/DomainUtils)](https://goreportcard.com/report/github.com/RouteHub-Link/DomainUtils)

<img width="1490" alt="image" src="https://github.com/RouteHub-Link/DomainUtils/assets/16222645/26482c9c-e0cd-4c71-825d-772d02b96b67">


## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Contributing](../CONTRIBUTING.md)

## About <a name = "about"></a>

This project is a Domain utility for validating an site. features are;

1. URL Validation: The service will validate a URL with Configuration options to allow or disallow certain URL patterns.
2. Seek The URL: The service will seek the URL and return the response status code. With some additional configuration options for more advanced use cases. (Does not does url validation please validate it first then create a seek task.)
3. DNS TXT Lookup: The service will perform a DNS TXT lookup on the domain of the URL and return the TXT record.

The service has asynq for queueing the these tasks and a REST API for the client to interact with the service.
For Using the service, the client will send a POST request to the service with the URL and the task to be performed.

The service will then queue the task and return a task ID to the client. The client can then use the task ID to check the status of the task and get the result of the task.

Other way to see task's is using asynqmon, a web-based monitoring tool for asynq.

## Getting Started <a name = "getting_started"></a>

### Prerequisites

What things you need to install the software and how to install them.

- go 1.22.1
- redis

### Service Dependencies

- [asynq](https://github.com/hibiken/asynq)
- [asynqmon](https://github.com/hibiken/asynqmon)
- [validator](https://github.com/RouteHub-Link/DomainUtils/tree/main/validator)
- [echo](https://github.com/labstack/echo)
- [koanf](https://github.com/knadh/koanf/)
- [mergo](https://github.com/darccio/mergo)

### Validator Dependencies

- [colly](https://github.com/gocolly/colly)
- [dns](https://github.com/miekg/dns)

### Installing as Service

Redis is a hard requirement for the service to run. Make sure you have redis installed and running.
If you wanna change the redis configuration, you can do it from main.go file.
You must run this following command to use as a service. (I will create a docker option for hosting the service soon.)

1. Clone the repository.

```bash
git clone https://github.com/RouteHub-Link/DomainUtils.git
cd DomainUtils
```

2. Build the service.

```bash
go build .
./DomainUtils --help
```

### 3. Starting the service.

3.1 Stating via Docker

```bash
docker-compose up --build
```

3.2 Starting the service.

Basicliy, you need 3 terminals to run the service. One for handling tasks, one for monitoring tasks and one for handling requests.

1. Starting for handling tasks.

```bash
    ./DomainUtils -s 1
```

2. Starting for monitoring tasks. (Runs asynqmon optional)

```bash
    ./DomainUtils -s 2
```

3. Starting for handling requests.

```bash
    ./DomainUtils
```

Endpoints and Responses

- POST /validate/url
- POST /validate/dns
  - Request
```json
    {
        "url": "https://www.google.com"
    }
```

  - Response
```json
    {
        "task_id": "task_id"
    }
```

- GET /validate/url/task_id
- GET /validate/dns/task_id
  - Response
```json
    {
        "ID": "732c5e2c-4dec-429a-9613-b1fe6427232b",
        "Queue": "url-validation",
        "Type": "url:validate",
        "Payload": "base64=",
        "State": 6,
        "MaxRetry": 10,
        "Retried": 0,
        "LastErr": "",
        "LastFailedAt": "0001-01-01T00:00:00Z",
        "Timeout": 120000000000,
        "Deadline": "2024-05-15T14:17:40+03:00",
        "Group": "",
        "NextProcessAt": "0001-01-01T00:00:00Z",
        "IsOrphaned": false,
        "Retention": 864000000000000,
        "CompletedAt": "2024-05-14T14:17:41+03:00",
        "Result": "base64="
    }
```

### Installing as Validator Package Only

```bash
- go get github.com/RouteHub-Link/DomainUtils/validator
```

Note check validator.go for configuration implementation.

```go
    var _validator = validator.DefaultValidator()
    isValid, err := _validator.ValidateURL(payload.Link)
    if err != nil {
        log.Println(err)
    }

    if isValid {
        log.Println("URL is valid")
    } else {
        log.Println("URL is not valid")
    }

```

```go
    customConfig := validator.CheckConfig{
        MaxRedirects:          2,       // Set your desired default values
        MaxSize:               4194304, // 4MB
        MaxURLLength:          2048,    // 2048 characters
        CheckForFile:          true,
        CheckIsReachable:      true,
        CannotEndWithSlash:    true,
        HTTPSRequired:         true,
        HTTPClientTimeout:     10 * time.Second,
        ContentTypeMustBeHTML: true,
    }

    _validator := validator.NewValidator(customConfig)
    isValid, err := _validator.ValidateURL(payload.Link)

    if isValid {
        isSiteValid, err = _validator.ValidateSite(link)
        // do something with bool & err
    }
```

## Usage <a name = "usage"></a>

For starting the service, you can use the following command:

```bash
go run .
```

Starting via config file please edit config.yaml file.
If there is not a config.yaml file, the service will use the default configuration and you will see the following output:

```bash
error loading config: open config.yaml: no such file or directory
```

Also you can use the following command to start the service.

```bash
go build .
./DomainUtils -r "127.0.0.1:6379" -p 1235
```

for more information about the flags, you can use the following command:

```bash
./DomainUtils -h
```

- For changin the serving mode, you can use the following command:

```bash
    -s, --serving-mode int8   Serves application as selected mode.
        -s 0 or --serving-mode 0
        0 : TaskReceiver
        1 : TaskServer
        2 : TaskMonitoring
        Default : 0

    ./DomainUtils -s 1
```
