package main

import (
	"bytes"
	"fmt"
	"regexp"
)

type ServiceHandler interface {
	Connect(service Service)
	Disconnect(service Service)
	Start(service Service) int
	Status(service Service) int
	Stop(service Service) int
}

type LinuxServiceHandler struct {
	handler ProtocolHandler
}

func (r *LinuxServiceHandler) Connect(service Service) {
	fmt.Println("connecting to server")
	r.handler.OpenConnection(service)
}

func (r *LinuxServiceHandler) Disconnect(service Service) {
	fmt.Println("disconnecting from server")
	r.handler.CloseConnection()
}

func (r *LinuxServiceHandler) Start(service Service) int {
	fmt.Println("starting service")

	service.action = "start"
	r.RunAction(service)

	return r.Status(service)
}

func (r *LinuxServiceHandler) Status(service Service) int {

	fmt.Println("determining service status")

	service.action = "status"
	status := ServiceStatusUnknown

	text := r.RunAction(service)

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

func (r *LinuxServiceHandler) RunAction(service Service) string {

	var buffer bytes.Buffer

	if service.sudo != "" {
		buffer.WriteString(fmt.Sprintf("echo '%s' | sudo -S service %s %s", service.sudo, service.name, service.action))
	} else {
		buffer.WriteString(fmt.Sprintf("sudo service %s %s", service.name, service.action))
	}

	text := r.handler.Run(buffer.String())

	return text
}

func (r LinuxServiceHandler) Stop(service Service) int {
	fmt.Println("stopping service")
	service.action = "stop"

	r.RunAction(service)

	return r.Status(service)
}
