package main

import (
	"fmt"
	"testing"
)

// Service is running with the pid 7112, with sudo
func TestLinuxSSHServiceHandlerStatus01(t *testing.T) {

	// given
	mock := MockSSHProtocolHandler{results: [10]string{"7112", "7112 ? 00:00:00 daemon"}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(LinuxSSHServiceHandler{handler: handler})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status",
		sudo:     "mysudo"}

	// when
	r.Connect(service)
	defer r.Disconnect(service)

	result := r.Status(service)

	// then
	if result != ServiceStatusStarted {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 2 {
		t.Error("Expected runs of 2, got ", mock.run)
	}

	if mock.runs[0] != "echo 'mysudo' | sudo -S find /var/run/ -name 'myname.pid' -exec cat {} \\; 2> /dev/null" {
		t.Error("Expected other, got ", mock.runs[0])
	}

	if mock.runs[1] != "ps -p 7112" {
		t.Error("Expected other, got ", mock.runs[1])
	}
}

// Service is running with the pid 7112, WITHOUT sudo
func TestLinuxSSHServiceHandlerStatus02(t *testing.T) {

	// given
	mock := MockSSHProtocolHandler{results: [10]string{"7112", "7112 ? 00:00:00 daemon"}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(LinuxSSHServiceHandler{handler: handler})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	r.Connect(service)
	defer r.Disconnect(service)

	result := r.Status(service)

	// then
	if result != ServiceStatusStarted {
		t.Error("Expected service started, got ", result)
	}

	if mock.run != 2 {
		t.Error("Expected runs of 2, got ", mock.run)
	}

	if mock.runs[0] != "find /var/run/ -name 'myname.pid' -exec cat {} \\; 2> /dev/null" {
		t.Error("Expected other, got ", mock.runs[0])
	}

	if mock.runs[1] != "ps -p 7112" {
		t.Error("Expected other, got ", mock.runs[1])
	}
}

// Status -> no permission to .pid file
func TestLinuxSSHServiceHandlerStatus03(t *testing.T) {

	// given
	mock := MockSSHProtocolHandler{results: [10]string{""}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(LinuxSSHServiceHandler{handler: handler})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	r.Connect(service)
	defer r.Disconnect(service)

	result := r.Status(service)

	// then
	if result != ServiceStatusUnknown {
		t.Error("Expected service unknown, got ", result)
	}

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "find /var/run/ -name 'myname.pid' -exec cat {} \\; 2> /dev/null" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

type MockSSHProtocolHandler struct {
	runs    [10]string
	results [10]string
	run     int
}

func (r *MockSSHProtocolHandler) OpenConnection(service Service) {
	fmt.Println("Mock Open connection")
}

func (r *MockSSHProtocolHandler) Run(cmd string) string {

	fmt.Println("sending cmd: ", cmd)

	r.runs[r.run] = cmd
	s := r.results[r.run]

	r.run += 1

	fmt.Println("got response ", s)

	return s
}

func (r *MockSSHProtocolHandler) CloseConnection() {
	fmt.Println("mock close connection")
}
