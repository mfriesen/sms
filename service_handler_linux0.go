package main

import (
	"bytes"
	"fmt"
	"regexp"
)

type LinuxServiceHandler struct {
}

func (r *LinuxServiceHandler) Start(service Service, handler ProtocolHandler) int {
	log.Info("starting %s service", service.name)

	service.action = "start"
	r.RunAction(service, handler)

	return r.Status(service, handler)
}

func (r *LinuxServiceHandler) Status(service Service, handler ProtocolHandler) int {

	log.Info("determining service %s status", service.name)

	service.action = "status"
	status := ServiceStatusUnknown

	text := r.RunAction(service, handler)

	if len(text) > 0 {

		rp := regexp.MustCompile("[0-9]+")

		if rp.MatchString(text) {
			status = ServiceStatusStarted
		} else {
			status = ServiceStatusStopped
		}
	}

	return status
}

func (r *LinuxServiceHandler) RunAction(service Service, handler ProtocolHandler) string {

	var buffer bytes.Buffer

	if service.sudo != "" {
		buffer.WriteString(fmt.Sprintf("echo '%s' | sudo -S service %s %s", service.sudo, service.name, service.action))
	} else {
		buffer.WriteString(fmt.Sprintf("sudo service %s %s", service.name, service.action))
	}

	text := handler.Run(buffer.String())

	return text
}

func (r *LinuxServiceHandler) Stop(service Service, handler ProtocolHandler) int {
	log.Info("stopping %s service", service.name)

	service.action = "stop"

	r.RunAction(service, handler)

	return r.Status(service, handler)
}

func (r LinuxServiceHandler) IsSupported() bool {
	return true
}
