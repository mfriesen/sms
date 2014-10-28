package main

import (
	"testing"
)

// test no parameters entered
func TestUsage01(t *testing.T) {
	// given
	vargs := []string{}

	// when
	_, err := usage(vargs, false)

	// then
	if err == nil {
		t.Error("Expected Errors, got none")
	}
}

// test correct STOP parameters entered
func TestUsage02(t *testing.T) {
	// given
	vargs := []string{"testuser@testhost", "servicename", "stop"}

	// when
	service, err := usage(vargs, false)

	// then
	if err != nil {
		t.Error("Expected NO Errors, got some")
	}

	if service.user != "testuser" {
		t.Error("Expected <user> testuser, got ", service.user)
	}

	if service.host != "testhost" {
		t.Error("Expected <host> testhost, got ", service.host)
	}

	if service.port != "22" {
		t.Error("Expected <port> 22, got ", service.port)
	}

	if service.action != "stop" {
		t.Error("Expected <action> stop, got ", service.action)
	}

	if service.name != "servicename" {
		t.Error("Expected servicename, got ", service.name)
	}
}

// test correct STATUS parameters entered
func TestUsage03(t *testing.T) {
	// given
	vargs := []string{"testuser@testhost", "servicename", "status"}

	// when
	service, err := usage(vargs, false)

	// then
	if err != nil {
		t.Error("Expected NO Errors, got some")
	}

	if service.user != "testuser" {
		t.Error("Expected <user> testuser, got ", service.user)
	}

	if service.host != "testhost" {
		t.Error("Expected <host> testhost, got ", service.host)
	}

	if service.port != "22" {
		t.Error("Expected <port> 22, got ", service.port)
	}

	if service.action != "status" {
		t.Error("Expected <action> status, got ", service.action)
	}

	if service.name != "servicename" {
		t.Error("Expected servicename, got ", service.name)
	}
}

// test correct START parameters entered
func TestUsage04(t *testing.T) {
	// given
	vargs := []string{"testuser@testhost", "servicename", "start"}

	// when
	service, err := usage(vargs, false)

	// then
	if err != nil {
		t.Error("Expected NO Errors, got some")
	}

	if service.user != "testuser" {
		t.Error("Expected <user> testuser, got ", service.user)
	}

	if service.host != "testhost" {
		t.Error("Expected <host> testhost, got ", service.host)
	}

	if service.port != "22" {
		t.Error("Expected <port> 22, got ", service.port)
	}

	if service.action != "start" {
		t.Error("Expected <action> start, got ", service.action)
	}

	if service.name != "servicename" {
		t.Error("Expected servicename, got ", service.name)
	}
}

// test correct parameters entered with port
func TestUsage05(t *testing.T) {
	// given
	vargs := []string{"testuser@testhost:25", "service", "stop"}

	// when
	service, err := usage(vargs, false)

	// then
	if err != nil {
		t.Error("Expected NO Errors, got some")
	}

	if service.user != "testuser" {
		t.Error("Expected <user> testuser, got ", service.user)
	}

	if service.host != "testhost" {
		t.Error("Expected <host> testhost, got ", service.host)
	}

	if service.port != "25" {
		t.Error("Expected <port> 25, got ", service.port)
	}
}

// test missing host
func TestUsage06(t *testing.T) {
	// given
	vargs := []string{"testuser:25"}

	// when
	_, err := usage(vargs, false)

	// then
	if err == nil {
		t.Error("Expected Errors, got none")
	}
}

// test --help
func TestUsage07(t *testing.T) {
	// given
	vargs := []string{"--help"}

	// when
	_, err := usage(vargs, false)

	// then
	if err != nil {
		t.Error("Expected NO Errors, got some")
	}
}

// test -v
func TestUsage08(t *testing.T) {
	// given
	vargs := []string{"-v", "testuser@testhost", "service", "stop"}

	// when
	_, err := usage(vargs, false)

	// then
	if err != nil {
		t.Error("Expected NO Errors, got some", err)
	}

	if log.IsDebug() == false {
		t.Error("Expected DEBUG")
	}
}

// test --verbose
func TestUsage09(t *testing.T) {
	// given
	vargs := []string{"--verbose", "testuser@testhost", "service", "stop"}

	// when
	_, err := usage(vargs, false)

	// then
	if err != nil {
		t.Error("Expected NO Errors, got some", err)
	}

	if log.IsDebug() == false {
		t.Error("Expected DEBUG")
	}
}
