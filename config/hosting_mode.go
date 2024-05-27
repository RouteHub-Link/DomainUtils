package config

import (
	"fmt"
	"strings"
)

type HostingMode int8

const (
	TaskClient HostingMode = iota
	TaskServer
	TaskMonitoring
	Debug
)

var HostingModeNames = map[HostingMode]string{
	TaskClient:     "TaskClient",
	TaskServer:     "TaskServer",
	TaskMonitoring: "TaskMonitoring",
	Debug:          "Debug",
}

var HostingModeIds = map[HostingMode]int8{
	TaskClient:     0,
	TaskServer:     1,
	TaskMonitoring: 2,
	Debug:          3,
}

func (a HostingMode) String() string {
	return HostingModeNames[a]
}

func (a HostingMode) Id() int8 {
	return HostingModeIds[a]
}

func ParseHostingModeFromString(s string) (HostingMode, error) {
	for k, v := range HostingModeNames {
		if strings.EqualFold(s, v) {
			return k, nil
		}
	}
	return 0, fmt.Errorf("invalid HostingMode: %s", s)
}

func ParseHostingModeFromInt(id int8) (HostingMode, error) {
	for k, v := range HostingModeIds {
		if id == v {
			return k, nil
		}
	}
	return 0, fmt.Errorf("invalid HostingMode: %d", id)

}

func (a *HostingMode) Set(value HostingMode) {
	*a = value
}

func (a HostingMode) Get() HostingMode {
	return a
}
