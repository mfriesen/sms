package main

import (
	"testing"
)

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

	r := ServiceHandler(&LinuxToWindowsServiceHandler{})
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

	r := ServiceHandler(&LinuxToWindowsServiceHandler{})
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

	r := ServiceHandler(&LinuxToWindowsServiceHandler{})
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

	r := ServiceHandler(&LinuxToWindowsServiceHandler{})
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

	r := ServiceHandler(&LinuxToWindowsServiceHandler{})
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
