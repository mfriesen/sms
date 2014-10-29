package main

import (
	"fmt"
	"strings"
)

type WindowsToWindowsServiceHandler struct {
}

func (r *WindowsToWindowsServiceHandler) Start(service Service, handler ProtocolHandler) int {

	cmd := fmt.Sprintf("sc \\\\%s start %s", service.host, service.name)
	handler.Run(cmd)

	return r.Status(service, handler)
}

func (r *WindowsToWindowsServiceHandler) Status(service Service, handler ProtocolHandler) int {

	status := ServiceStatusUnknown
	cmd := fmt.Sprintf("sc \\\\%s query %s", service.host, service.name)

	text := handler.Run(cmd)

	if strings.Contains(text, "RUNNING") {
		status = ServiceStatusStarted
	} else if strings.Contains(text, "STOPPED") {
		status = ServiceStatusStopped
	}

	return status
}

func (r WindowsToWindowsServiceHandler) Stop(service Service, handler ProtocolHandler) int {

	cmd := fmt.Sprintf("sc \\\\%s stop %s", service.host, service.name)
	handler.Run(cmd)

	return r.Status(service, handler)
}
