package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type ServiceHandler interface {
	Connect(service Service)
	Disconnect(service Service)
	Start(service Service)
	Status(service Service) int
	Stop(service Service)
}

type LinuxSSHServiceHandler struct {
	handler ProtocolHandler
}

func (r LinuxSSHServiceHandler) Connect(service Service) {
	fmt.Println("connecting to server")
	r.handler.OpenConnection(service)
}

func (r LinuxSSHServiceHandler) Disconnect(service Service) {
	fmt.Println("disconnecting from server")
	r.handler.CloseConnection()
}

func (r LinuxSSHServiceHandler) Start(service Service) {
	fmt.Println("starting service")
}

// find /var/run/ -name 'jenkins.pid' -exec cat {} \; 2> /dev/null | xargs ps -p
func (r LinuxSSHServiceHandler) Status(service Service) int {

	status := ServiceStatusUnknown

	fmt.Println("determining service status")

	var buffer bytes.Buffer

	if service.sudo != "" {
		buffer.WriteString(fmt.Sprintf("echo '%s' | sudo -S ", service.sudo))
	}

	buffer.WriteString(fmt.Sprintf("find /var/run/ -name '%s.pid' -exec cat {} \\; 2> /dev/null", service.name))

	pid := r.handler.Run(buffer.String())

	if _, err := strconv.Atoi(pid); err == nil {

		result := r.handler.Run(fmt.Sprintf("ps -p %s", pid))

		if strings.Contains(result, pid) {
			status = ServiceStatusStarted
		}
	}

	return status
}

func (r LinuxSSHServiceHandler) Stop(service Service) {
	fmt.Println("stopping service")
}
