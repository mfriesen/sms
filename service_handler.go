package main

import (
	"bytes"
	"fmt"
	//"github.com/jcelliott/lumber"
	"regexp"
	"strings"
)

type ServiceHandler interface {
	Start(service Service, handler ProtocolHandler) int
	Status(service Service, handler ProtocolHandler) int
	Stop(service Service, handler ProtocolHandler) int
	IsSupported(handler ProtocolHandler) bool
}

type ServiceExecServiceHandler struct {
}

func (r *ServiceExecServiceHandler) Start(service Service, protocol ProtocolHandler) int {
	log.Info("starting %s service", service.name)

	service.action = "start"
	r.RunAction(service, protocol)

	return r.Status(service, protocol)
}

func (r *ServiceExecServiceHandler) Status(service Service, protocol ProtocolHandler) int {

	log.Info("determining service %s status", service.name)

	service.action = "status"
	status := ServiceStatusUnknown

	text := r.RunAction(service, protocol)

	if len(text) > 0 {

		rp0 := regexp.MustCompile("( start)|( running)")
		rp1 := regexp.MustCompile("( stop)")

		if rp0.MatchString(text) {
			status = ServiceStatusStarted
		} else if rp1.MatchString(text) {
			status = ServiceStatusStopped
		}
	}

	return status
}

func (r *ServiceExecServiceHandler) RunAction(service Service, protocol ProtocolHandler) string {

	var buffer bytes.Buffer

	if service.sudo != "" {
		buffer.WriteString(fmt.Sprintf("echo '%s' | sudo -S service %s %s", service.sudo, service.name, service.action))
	} else {
		buffer.WriteString(fmt.Sprintf("sudo service %s %s", service.name, service.action))
	}

	stdout, err := protocol.Run(service, buffer.String())

	if err != nil && strings.Contains(err.Error(), "sudo:") {
		log.Fatal("'--sudo' parameter required for this service")
	}

	return stdout
}

func (r *ServiceExecServiceHandler) Stop(service Service, protocol ProtocolHandler) int {
	log.Info("stopping %s service", service.name)

	service.action = "stop"

	r.RunAction(service, protocol)

	return r.Status(service, protocol)
}

func (r *ServiceExecServiceHandler) IsSupported(protocol ProtocolHandler) bool {
	return isCommandSupported(protocol, "service")
}

func isCommandSupported(protocol ProtocolHandler, cmd string) bool {

	log.Debug("looking for executable '%s'", cmd)

	_, err := protocol.Run(Service{}, cmd)

	return err == nil
}

func checkCommandSupported(stdout string, stderr string) bool {
	return !(strings.Contains(stdout, "not found") || strings.Contains(stderr, "not found") ||
		strings.Contains(stdout, "not recognized") || strings.Contains(stderr, "not recognized") ||
		strings.Contains(stdout, "not exist") || strings.Contains(stderr, "not exist"))
}

type SambaServiceHandler struct {
}

func (r *SambaServiceHandler) Start(service Service, protocol ProtocolHandler) int {
	cmd := fmt.Sprintf("net rpc service start %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)
	protocol.Run(service, cmd)

	return r.Status(service, protocol)
}

func (r *SambaServiceHandler) Status(service Service, protocol ProtocolHandler) int {

	status := ServiceStatusUnknown

	cmd := fmt.Sprintf("net rpc service status %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)
	stdout, _ := protocol.Run(service, cmd)

	if strings.Contains(stdout, "is running") {
		status = ServiceStatusStarted
	} else if strings.Contains(stdout, "is stopped") {
		status = ServiceStatusStopped
	}

	return status
}

func (r *SambaServiceHandler) Stop(service Service, protocol ProtocolHandler) int {

	cmd := fmt.Sprintf("net rpc service stop %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)
	protocol.Run(service, cmd)

	return r.Status(service, protocol)
}

func (r *SambaServiceHandler) IsSupported(protocol ProtocolHandler) bool {
	return isCommandSupported(protocol, "net")
}

type ScExecServiceHandler struct {
}

func (r *ScExecServiceHandler) Start(service Service, protocol ProtocolHandler) int {

	cmd := fmt.Sprintf("sc \\\\%s start %s", service.host, service.name)
	protocol.Run(service, cmd)

	return r.Status(service, protocol)
}

func (r *ScExecServiceHandler) Status(service Service, protocol ProtocolHandler) int {

	status := ServiceStatusUnknown
	cmd := fmt.Sprintf("sc \\\\%s query %s", service.host, service.name)

	stdout, _ := protocol.Run(service, cmd)

	if strings.Contains(stdout, "RUNNING") {
		status = ServiceStatusStarted
	} else if strings.Contains(stdout, "STOPPED") {
		status = ServiceStatusStopped
	}

	return status
}

func (r *ScExecServiceHandler) Stop(service Service, protocol ProtocolHandler) int {

	cmd := fmt.Sprintf("sc \\\\%s stop %s", service.host, service.name)
	protocol.Run(service, cmd)

	return r.Status(service, protocol)
}

func (r *ScExecServiceHandler) IsSupported(protocol ProtocolHandler) bool {
	return isCommandSupported(protocol, "sc query")
}
