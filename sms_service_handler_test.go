package main

import (
	"testing"
)

// Service is running with the pid 7112
func TestLinuxServiceHandlerStatus01(t *testing.T) {

	// given
	mock := MockProtocolHandler{results: [10]string{"Service is running with the pid 7112"}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&LinuxServiceHandler{handler: handler})
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

	if mock.run != 1 {
		t.Error("Expected runs of 1, got ", mock.run)
	}

	if mock.runs[0] != "sudo service myname status" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// service is running (17687), with SUDO
func TestLinuxServiceHandlerStatus02(t *testing.T) {

	// given
	mock := MockProtocolHandler{results: [10]string{"service is running (17687)"}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&LinuxServiceHandler{handler: handler})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		sudo:     "mysudo",
		action:   "status"}

	// when
	r.Connect(service)
	defer r.Disconnect(service)

	result := r.Status(service)

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
func TestLinuxServiceHandlerStatus03(t *testing.T) {

	// given
	mock := MockProtocolHandler{results: [10]string{""}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&LinuxServiceHandler{handler: handler})
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

	if mock.runs[0] != "sudo service myname status" {
		t.Error("Expected other, got ", mock.runs[0])
	}
}

// Started Service successful
func TestLinuxServiceHandlerStart01(t *testing.T) {

	// given
	mock := MockProtocolHandler{results: [10]string{"Started service with pid 7112", "Service is running with the pid 7112"}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&LinuxServiceHandler{handler: handler})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	r.Connect(service)
	defer r.Disconnect(service)

	result := r.Start(service)

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
func TestLinuxServiceHandlerStop01(t *testing.T) {

	// given
	mock := MockProtocolHandler{results: [10]string{"Stopped service", "Service is stopped"}}

	handler := ProtocolHandler(&mock)

	r := ServiceHandler(&LinuxServiceHandler{handler: handler})
	service := Service{
		user:     "myuser",
		password: "mypass",
		host:     "myhost",
		name:     "myname",
		action:   "status"}

	// when
	r.Connect(service)
	defer r.Disconnect(service)

	result := r.Stop(service)

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

func (r *MockProtocolHandler) Run(cmd string) string {

	log.Info("mock sending cmd: ", cmd)

	r.runs[r.run] = cmd
	s := r.results[r.run]

	r.run += 1

	log.Info("mock got response ", s)

	return s
}

func (r *MockProtocolHandler) CloseConnection() {
	log.Info("mock close connection")
}
