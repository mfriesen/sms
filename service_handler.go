package main

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type ServiceHandler interface {
	Start(service Service, protocol ProtocolHandler) (int, error)
	Status(service Service, protocol ProtocolHandler) (int, error)
	Stop(service Service, protocol ProtocolHandler) (int, error)
	Search(service Service, protocol ProtocolHandler) ([]string, error)
	IsSupported(protocol ProtocolHandler) bool
}

type ServiceExecServiceHandler struct {
}

func (r *ServiceExecServiceHandler) Search(service Service, protocol ProtocolHandler) ([]string, error) {
	log.Info("search for %s service", service.name)

	cmd := "service --status-all"
	stdout, err := protocol.Run(service, cmd)

	return Search(service, stdout, err)
}

func (r *ServiceExecServiceHandler) Start(service Service, protocol ProtocolHandler) (int, error) {
	log.Info("starting %s service", service.name)
	cmd := fmt.Sprintf("service %s start", service.name)
	cmd = r.AddSudo(cmd, service)
	return StartOrStopWitbRetry(service, protocol, r, cmd, ServiceStatusStarted)
}

func (r *ServiceExecServiceHandler) Status(service Service, protocol ProtocolHandler) (int, error) {

	log.Info("determining service %s status", service.name)

	status := ServiceStatusUnknown
	cmd := fmt.Sprintf("service %s status", service.name)
	cmd = r.AddSudo(cmd, service)

	stdout, err := protocol.Run(service, cmd)

	if len(stdout) > 0 {

		rp0 := regexp.MustCompile("( start)|( is running)")
		rp1 := regexp.MustCompile("( stop)|( is not running)")

		if rp0.MatchString(stdout) {
			status = ServiceStatusStarted
		} else if rp1.MatchString(stdout) {
			status = ServiceStatusStopped
		}
	}

	return status, err
}

func (r *ServiceExecServiceHandler) AddSudo(cmd string, service Service) string {

	if service.sudo != "" {
		return fmt.Sprintf("echo '%s' | sudo -S %s", service.sudo, cmd)
	} else {
		return fmt.Sprintf("sudo %s", cmd)
	}
}

func (r *ServiceExecServiceHandler) Stop(service Service, protocol ProtocolHandler) (int, error) {
	log.Info("stopping %s service", service.name)

	cmd := fmt.Sprintf("service %s stop", service.name)
	cmd = r.AddSudo(cmd, service)
	return StartOrStopWitbRetry(service, protocol, r, cmd, ServiceStatusStopped)
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

func (r *SambaServiceHandler) Search(service Service, protocol ProtocolHandler) ([]string, error) {
	log.Info("search for %s service", service.name)
	cmd := fmt.Sprintf("net rpc service list -I %s -U %s%%%s", service.host, service.user, service.password)

	stdout, err := protocol.Run(service, cmd)

	return Search(service, stdout, err)
}

func (r *SambaServiceHandler) Start(service Service, protocol ProtocolHandler) (int, error) {
	cmd := fmt.Sprintf("net rpc service start %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)
	return StartOrStopWitbRetry(service, protocol, r, cmd, ServiceStatusStarted)
}

func (r *SambaServiceHandler) Stop(service Service, protocol ProtocolHandler) (int, error) {
	cmd := fmt.Sprintf("net rpc service stop %s -I %s -U %s%%%s", service.name, service.host, service.user, service.password)
	return StartOrStopWitbRetry(service, protocol, r, cmd, ServiceStatusStopped)
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

func (r *SambaServiceHandler) IsSupported(protocol ProtocolHandler) bool {
	return strings.Contains(runtime.GOOS, "linux")
}

type ScExecServiceHandler struct {
	errorCount int
}

func (r *ScExecServiceHandler) Search(service Service, protocol ProtocolHandler) ([]string, error) {
	log.Info("search for %s service", service.name)

	cmd := fmt.Sprintf("wmic /node:'%s' service where (name like '%%%s%%') get name", service.host, service.name)

	stdout, err := protocol.Run(service, cmd)
	return Search(service, stdout, err)
}

func (r *ScExecServiceHandler) Start(service Service, protocol ProtocolHandler) (int, error) {
	cmd := fmt.Sprintf("sc \\\\%s start %s", service.host, service.name)
	return StartOrStopWitbRetry(service, protocol, r, cmd, ServiceStatusStarted)
}

func (r *ScExecServiceHandler) Stop(service Service, protocol ProtocolHandler) (int, error) {
	cmd := fmt.Sprintf("sc \\\\%s stop %s", service.host, service.name)
	return StartOrStopWitbRetry(service, protocol, r, cmd, ServiceStatusStopped)
}

func (r *ScExecServiceHandler) Status(service Service, protocol ProtocolHandler) (int, error) {

	status := ServiceStatusUnknown
	cmd := fmt.Sprintf("sc \\\\%s query %s", service.host, service.name)

	stdout, err := protocol.Run(service, cmd)

	// windows returns right away, give it some time to update the service's status
	time.Sleep(time.Duration(1000) * time.Millisecond)

	if strings.Contains(stdout, "_PENDING") {
		if r.errorCount < 60 {
			fmt.Print(".")
			time.Sleep(time.Duration(500) * time.Millisecond)
			r.errorCount++

			status, err = r.Status(service, protocol)
		}
	} else if strings.Contains(stdout, "RUNNING") {
		status = ServiceStatusStarted
	} else if strings.Contains(stdout, "STOPPED") {
		status = ServiceStatusStopped
	}

	return status, err
}

func (r *ScExecServiceHandler) IsSupported(protocol ProtocolHandler) bool {
	return strings.Contains(runtime.GOOS, "windows")
}

func StartOrStopWitbRetry(service Service, protocol ProtocolHandler, serviceHandler ServiceHandler, cmd string, wantedStatus int) (int, error) {

	var err error
	status := ServiceStatusUnknown

	_, retErr := protocol.Run(service, cmd)

	i := 0
	for status != wantedStatus {

		status, err = serviceHandler.Status(service, protocol)

		// set retErr to error from status only if it's never been set
		// or call to Status returned no error
		if retErr == nil || err == nil {
			retErr = err
		}

		if i == 30 || retErr != nil {
			status = ServiceStatusUnknown
			break
		} else {
			fmt.Print(".")
			time.Sleep(time.Duration(1000) * time.Millisecond)
			i++
		}
	}

	return status, retErr
}

func Search(service Service, stdout string, err error) ([]string, error) {

	list := []string{}

	if err == nil {

		rp := regexp.MustCompile(fmt.Sprintf(".*%s.*", service.name))

		strs := strings.Split(stdout, "\n")

		for _, element := range strs {

			if rp.MatchString(element) {
				list = append(list, strings.Trim(element, " "))
			}
		}
	}

	return list, err

}
