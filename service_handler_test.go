package main

import (
	"testing"
)

// Service is running with the pid 7112
func TestServiceExecServiceHandlerStatus01(t *testing.T) {

	// given
	mock := MockProtocolHandler{results: [10]string{"Service is running with the pid 7112"}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&ServiceExecServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	//r.Connect(service)
	//defer r.Disconnect(service)

	result := r.Status(service, handler)

	// then
	if result != ServiceStatusStarted {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "sudo service myname status" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// service is running (17687), with SUDO
func TestServiceExecServiceHandlerStatus02(t *testing.T) {

	// given
	mock := MockProtocolHandler{results: [10]string{"service is running (17687)"}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&ServiceExecServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		sudo:     "mysudo",
		action:   "status"}

	// when
	//r.Connect(service)
	//defer r.Disconnect(service)

	result := r.Status(service, handler)

	// then
	if result != ServiceStatusStarted {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "echo 'mysudo' | sudo -S service myname status" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// text return is empty
func TestServiceExecServiceHandlerStatus03(t *testing.T) {

	// given
	mock := MockProtocolHandler{results: [10]string{""}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&ServiceExecServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	//r.Connect(service)
	//defer r.Disconnect(service)

	result := r.Status(service, handler)

	// then
	if result != ServiceStatusUnknown {
		t.Error("Expected service unknown, got ", result)
	}

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "sudo service myname status" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// Started Service successful
func TestServiceExecServiceHandlerStart01(t *testing.T) {

	// given
	mock := MockProtocolHandler{results: [10]string{"Started service with pid 7112", "Service is running with the pid 7112"}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&ServiceExecServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	//r.Connect(service)
	//defer r.Disconnect(service)

	result := r.Start(service, handler)

	// then
	if result != ServiceStatusStarted {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 2 {
		t.Error("Expected runs of 2, got ", mock.run)
	}

	if mock.runs[0] != "sudo service myname start" {
		t.Error("Expected other, got ", mock.runs[0])
	}

	if mock.runs[1] != "sudo service myname status" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// Stopped Service successful
func TestServiceExecServiceHandlerStop01(t *testing.T) {

	// given
	mock := MockProtocolHandler{results: [10]string{"Stopped service", "Service is stopped"}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&ServiceExecServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	//r.Connect(service)
	//defer r.Disconnect(service)

	result := r.Stop(service, handler)

	// then
	if result != ServiceStatusStopped {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 2 {
		t.Error("Expected runs of 2, got ", mock.run)
	}

	if mock.runs[0] != "sudo service myname stop" {
		t.Error("Expected other, got ", mock.runs[0])
	}

	if mock.runs[1] != "sudo service myname status" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

type MockProtocolHandler struct {
	runs    [10]string
	results [10]string
	run     int
}

func (r *MockProtocolHandler) OpenConnection(service Service) {
	log.Info("Mock Open connection")
}

func (r *MockProtocolHandler) Run(service Service, cmd string) (string, string) {

	log.Info("mock sending cmd: ", cmd)

	r.runs[r.run] = cmd
	s := r.results[r.run]

	r.run += 1

	log.Info("mock got response ", s)

	return s, ""
}

func (r *MockProtocolHandler) CloseConnection(service Service) {
	log.Info("mock close connection")
}

func (r *MockProtocolHandler) IsSupported(service Service) bool {
	return true
}

func (r *MockProtocolHandler) IsPasswordNeeded(service Service) bool {
	return true
}

// Service is running
func TestLinuxToWindowsStatus01(t *testing.T) {
	// given
	mock := MockProtocolHandler{results: [10]string{
		`myname service is running.
Configuration details:
        Controls Accepted    = 0x45
        Service Type         = 0x10
        Start Type           = 0x2
        Error Control        = 0x1
        Tag ID               = 0x0
        Executable Path      = 
        Load Order Group     =
        Dependencies         = /
        Start Name           = myname
        Display Name         = myname`,
	}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&SambaServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	result := r.Status(service, handler)

	// then
	if result != ServiceStatusStarted {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "net rpc service status myname -I myhost -U myuser%mypass" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// Service is stopped
func TestLinuxToWindowsStatus02(t *testing.T) {
	// given
	mock := MockProtocolHandler{results: [10]string{
		`myname service is stopped.
Configuration details:
        Controls Accepted    = 0x0
        Service Type         = 0x10
        Start Type           = 0x3
        Error Control        = 0x1
        Tag ID               = 0x0
        Executable Path      =
        Load Order Group     =
        Dependencies         =
        Start Name           = LocalSystem
        Display Name         = myname
`,
	}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&SambaServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	result := r.Status(service, handler)

	// then
	if result != ServiceStatusStopped {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "net rpc service status myname -I myhost -U myuser%mypass" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// Service is unknown
func TestLinuxToWindowsStatus03(t *testing.T) {
	// given
	mock := MockProtocolHandler{results: [10]string{
		`Failed to open service.  [WERR_NO_SUCH_SERVICE]`,
	}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&SambaServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	result := r.Status(service, handler)

	// then
	if result != ServiceStatusUnknown {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "net rpc service status myname -I myhost -U myuser%mypass" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// Start Service
func TestLinuxToWindowsStart01(t *testing.T) {
	// given
	mock := MockProtocolHandler{results: [10]string{"",
		`myname service is running.
Configuration details:
        Controls Accepted    = 0x45
        Service Type         = 0x10
        Start Type           = 0x2
        Error Control        = 0x1
        Tag ID               = 0x0
        Executable Path      = 
        Load Order Group     =
        Dependencies         = /
        Start Name           = myname
        Display Name         = myname`,
	}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&SambaServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "start"}

	// when
	result := r.Start(service, handler)

	// then
	if result != ServiceStatusStarted {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 2 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "net rpc service start myname -I myhost -U myuser%mypass" {
		t.Error("Expected other, got ", mock.runs[0])
	}

	if mock.runs[1] != "net rpc service status myname -I myhost -U myuser%mypass" {
		t.Error("Expected other, got ", mock.runs[1])
	}
}

// Stop Service
func TestLinuxToWindowsStop01(t *testing.T) {
	// given
	mock := MockProtocolHandler{results: [10]string{"",
		`myname service is stopped.
Configuration details:
        Controls Accepted    = 0x0
        Service Type         = 0x10
        Start Type           = 0x3
        Error Control        = 0x1
        Tag ID               = 0x0
        Executable Path      =
        Load Order Group     =
        Dependencies         =
        Start Name           = LocalSystem
        Display Name         = myname
`,
	}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&SambaServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "start"}

	// when
	result := r.Stop(service, handler)

	// then
	if result != ServiceStatusStopped {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 2 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "net rpc service stop myname -I myhost -U myuser%mypass" {
		t.Error("Expected other, got ", mock.runs[0])
	}

	if mock.runs[1] != "net rpc service status myname -I myhost -U myuser%mypass" {
		t.Error("Expected other, got ", mock.runs[1])
	}
}

// Service is running
func TestWindowsToWindowsStatus01(t *testing.T) {
	// given
	mock := MockProtocolHandler{results: [10]string{
		`SERVICE_NAME: myname
	       TYPE               : 10
	       WIN32_OWN_PROCESS
	       STATE              : 4  RUNNING
	                             (STOPPABLE, NOT_PAUSABLE, ACCEPTS_SHUTDOWN)
	       WIN32_EXIT_CODE    : 0  (0x0)
	       SERVICE_EXIT_CODE  : 0  (0x0)
	       CHECKPOINT         : 0x0
	       WAIT_HINT          : 0x0`,
	}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&ScExecServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	result := r.Status(service, handler)

	// then
	if result != ServiceStatusStarted {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "sc \\\\myhost query myname" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// Service is stopped
func TestWindowsToWindowsStatus02(t *testing.T) {
	// given
	mock := MockProtocolHandler{results: [10]string{
		`SERVICE_NAME: myname
        TYPE               : 10  WIN32_OWN_PROCESS
        STATE              : 1  STOPPED
        WIN32_EXIT_CODE    : 1067  (0x42b)
        SERVICE_EXIT_CODE  : 0  (0x0)
        CHECKPOINT         : 0x0
        WAIT_HINT          : 0x0`,
	}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&ScExecServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	result := r.Status(service, handler)

	// then
	if result != ServiceStatusStopped {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "sc \\\\myhost query myname" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// Service is unknown
func TestWindowsToWindowsStatus03(t *testing.T) {
	// given
	mock := MockProtocolHandler{results: [10]string{
		`[SC] EnumQueryServicesStatus:OpenService FAILED 1060:

The specified service does not exist as an installed service.`,
	}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&ScExecServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	result := r.Status(service, handler)

	// then
	if result != ServiceStatusUnknown {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "sc \\\\myhost query myname" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// Start Service
func TestWindowsToWindowsStart01(t *testing.T) {
	// given
	mock := MockProtocolHandler{results: [10]string{"",
		`SERVICE_NAME: myname
	       TYPE               : 10
	       WIN32_OWN_PROCESS
	       STATE              : 4  RUNNING
	                             (STOPPABLE, NOT_PAUSABLE, ACCEPTS_SHUTDOWN)
	       WIN32_EXIT_CODE    : 0  (0x0)
	       SERVICE_EXIT_CODE  : 0  (0x0)
	       CHECKPOINT         : 0x0
	       WAIT_HINT          : 0x0`,
	}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&ScExecServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "start"}

	// when
	result := r.Start(service, handler)

	// then
	if result != ServiceStatusStarted {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 2 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "sc \\\\myhost start myname" {
		t.Error("Expected other, got ", mock.runs[0])
	}

	if mock.runs[1] != "sc \\\\myhost query myname" {
		t.Error("Expected other, got ", mock.runs[1])
	}
}

// Stop Service
func TestWindowsToWindowsStop01(t *testing.T) {
	// given
	mock := MockProtocolHandler{results: [10]string{"",
		`SERVICE_NAME: myname
        TYPE               : 10  WIN32_OWN_PROCESS
        STATE              : 1  STOPPED
        WIN32_EXIT_CODE    : 1067  (0x42b)
        SERVICE_EXIT_CODE  : 0  (0x0)
        CHECKPOINT         : 0x0
        WAIT_HINT          : 0x0`,
	}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&ScExecServiceHandler{})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "start"}

	// when
	result := r.Stop(service, handler)

	// then
	if result != ServiceStatusStopped {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 2 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "sc \\\\myhost stop myname" {
		t.Error("Expected other, got ", mock.runs[0])
	}

	if mock.runs[1] != "sc \\\\myhost query myname" {
		t.Error("Expected other, got ", mock.runs[1])
	}
}
