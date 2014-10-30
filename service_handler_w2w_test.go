package main

import (
	"testing"
)

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

	r := ServiceHandler(&WindowsToWindowsServiceHandler{})
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

	r := ServiceHandler(&WindowsToWindowsServiceHandler{})
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

	r := ServiceHandler(&WindowsToWindowsServiceHandler{})
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

	r := ServiceHandler(&WindowsToWindowsServiceHandler{})
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

	r := ServiceHandler(&WindowsToWindowsServiceHandler{})
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
