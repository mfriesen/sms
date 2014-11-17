package main

import (
	"bytes"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type ServiceHandler interface {
	Start(service Service, handler ProtocolHandler) (int, error)
	Status(service Service, handler ProtocolHandler) (int, error)
	Stop(service Service, handler ProtocolHandler) (int, error)
	IsSupported(handler ProtocolHandler) bool
}

type ServiceExecServiceHandler struct {
}

func (r *ServiceExecServiceHandler) Start(service Service, protocol ProtocolHandler) (int, error) {
	log.Info("starting %s service", service.name)

	status := ServiceStatusUnknown
	service.action = "start"
	_, err := r.RunAction(service, protocol)

	if err == nil {
		status, err = r.Status(service, protocol)
	}

	return status, err
}

func (r *ServiceExecServiceHandler) Status(service Service, protocol ProtocolHandler) (int, error) {

	log.Info("determining service %s status", service.name)

	service.action = "status"
	status := ServiceStatusUnknown

	text, err := r.RunAction(service, protocol)

	if len(text) > 0 {

		rp0 := regexp.MustCompile("( start)|( running)")
		rp1 := regexp.MustCompile("( stop)")

		if rp0.MatchString(text) {
			status = ServiceStatusStarted
		} else if rp1.MatchString(text) {
			status = ServiceStatusStopped
		}
	}

	return status, err
}

func (r *ServiceExecServiceHandler) RunAction(service Service, protocol ProtocolHandler) (string, error) {

	var buffer bytes.Buffer

	if service.sudo != "" {
		buffer.WriteString(fmt.Sprintf("echo '%s' | sudo -S service %s %s", service.sudo, service.name, service.action))
	} else {
		buffer.WriteString(fmt.Sprintf("sudo service %s %s", service.name, service.action))
	}

	stdout, err := protocol.Run(service, buffer.String())

	return stdout, err
}

func (r *ServiceExecServiceHandler) Stop(service Service, protocol ProtocolHandler) (int, error) {
	log.Info("stopping %s service", service.name)

	status := ServiceStatusUnknown
	service.action = "stop"

	_, err := r.RunAction(service, protocol)

	if err == nil {
		status, err = r.Status(service, protocol)
	}

	return status, err
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

func (r *SambaServiceHandler) Start(service Service, protocol ProtocolHandler) (int, error) {

	status := ServiceStatusUnknown
	cmd := fmt.Sprintf("net rpc service start %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)

	_, err := protocol.Run(service, cmd)

	if err == nil {
		status, err = r.Status(service, protocol)
	}

	return status, err
}

func (r *SambaServiceHandler) Status(service Service, protocol ProtocolHandler) (int, error) {

	status := ServiceStatusUnknown

	cmd := fmt.Sprintf("net rpc service status %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)
	stdout, err := protocol.Run(service, cmd)

	if strings.Contains(stdout, "is running") {
		status = ServiceStatusStarted
	} else if strings.Contains(stdout, "is stopped") {
		status = ServiceStatusStopped
	}

	return status, err
}

func (r *SambaServiceHandler) Stop(service Service, protocol ProtocolHandler) (int, error) {

	status := ServiceStatusUnknown
	cmd := fmt.Sprintf("net rpc service stop %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)

	_, err := protocol.Run(service, cmd)

	if err == nil {
		status, err = r.Status(service, protocol)
	}

	return status, err
}

func (r *SambaServiceHandler) IsSupported(protocol ProtocolHandler) bool {
	return strings.Contains(runtime.GOOS, "linux")
}

type ScExecServiceHandler struct {
	errorCount int
}

func (r *ScExecServiceHandler) Start(service Service, protocol ProtocolHandler) (int, error) {

	status := ServiceStatusUnknown
	cmd := fmt.Sprintf("sc \\\\%s start %s", service.host, service.name)
	_, err := protocol.Run(service, cmd)

	if err == nil {
		status, err = r.Status(service, protocol)
	}

	return status, err
}

func (r *ScExecServiceHandler) Status(service Service, protocol ProtocolHandler) (int, error) {

	status := ServiceStatusUnknown
	cmd := fmt.Sprintf("sc \\\\%s query %s", service.host, service.name)

	stdout, err := protocol.Run(service, cmd)

	// windows returns right away, give it some time to update the service's status
	time.Sleep(time.Duration(1000) * time.Millisecond)

	if strings.Contains(stdout, "_PENDING") {
		if r.errorCount < 60 {
			time.Sleep(time.Duration(500) * time.Millisecond)
			r.errorCount++
			fmt.Println(r.errorCount)

			status, err = r.Status(service, protocol)
		}
	} else if strings.Contains(stdout, "RUNNING") {
		status = ServiceStatusStarted
	} else if strings.Contains(stdout, "STOPPED") {
		status = ServiceStatusStopped
	}

	return status, err
}

func (r *ScExecServiceHandler) Stop(service Service, protocol ProtocolHandler) (int, error) {

	status := ServiceStatusUnknown
	cmd := fmt.Sprintf("sc \\\\%s stop %s", service.host, service.name)
	_, err := protocol.Run(service, cmd)

	if err == nil {
		status, err = r.Status(service, protocol)
	}

	return status, err
}

func (r *ScExecServiceHandler) IsSupported(protocol ProtocolHandler) bool {
	return strings.Contains(runtime.GOOS, "windows")
}
