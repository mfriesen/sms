package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
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

func (r *LinuxServiceHandler) IsSupported() bool {
	return runtime.GOOS == "linux"
}

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

func (r *LinuxToWindowsServiceHandler) Stop(service Service, handler ProtocolHandler) int {

	cmd := fmt.Sprintf("net rpc service stop %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)
	handler.Run(cmd)

	return r.Status(service, handler)
}

func (r *LinuxToWindowsServiceHandler) IsSupported() bool {
	return runtime.GOOS == "linux"
}

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

func (r *WindowsToWindowsServiceHandler) Stop(service Service, handler ProtocolHandler) int {

	cmd := fmt.Sprintf("sc \\\\%s stop %s", service.host, service.name)
	handler.Run(cmd)

	return r.Status(service, handler)
}

func (r *WindowsToWindowsServiceHandler) IsSupported() bool {

	if runtime.GOOS == "windows" {
		_, error := exec.LookPath("sc.exe")
		if error != nil {
			log.Fatal("cannot find sc.exe")
		}
	}

	return true
}
