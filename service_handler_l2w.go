package main

/*
net rpc service status SERVICE_NAME -I HOST -U USER%PASS

SERVICE_NAME service is running.
*/
import (
	"fmt"
	"strings"
)

type LinuxToWindowsServiceHandler struct {
}

func (r *LinuxToWindowsServiceHandler) Start(service Service, handler ProtocolHandler) int {
	cmd := fmt.Sprintf("net rpc service start %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)
	handler.Run(cmd)

	return r.Status(service, handler)
}

func (r *LinuxToWindowsServiceHandler) Status(service Service, handler ProtocolHandler) int {

	status := ServiceStatusUnknown

	cmd := fmt.Sprintf("net rpc service status %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)
	text := handler.Run(cmd)

	if strings.Contains(text, "is running") {
		status = ServiceStatusStarted
	} else if strings.Contains(text, "is stopped") {
		status = ServiceStatusStopped
	}

	return status
}

func (r LinuxToWindowsServiceHandler) Stop(service Service, handler ProtocolHandler) int {

	cmd := fmt.Sprintf("net rpc service stop %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)
	handler.Run(cmd)

	return r.Status(service, handler)
}

func (r LinuxToWindowsServiceHandler) IsSupported() bool {
	return true
}
