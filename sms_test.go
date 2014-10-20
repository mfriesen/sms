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

// test correct parameters entered
func TestUsage02(t *testing.T) {
	// given
	vargs := []string{"testuser@testhost"}

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
}

// test correct parameters entered with port
func TestUsage03(t *testing.T) {
	// given
	vargs := []string{"testuser@testhost:25"}

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
func TestUsage04(t *testing.T) {
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
func TestUsage05(t *testing.T) {
	// given
	vargs := []string{"--help"}

	// when
	_, err := usage(vargs, false)

	// then
	if err != nil {
		t.Error("Expected NO Errors, got some")
	}
}
