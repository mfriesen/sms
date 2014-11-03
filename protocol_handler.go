package main

import (
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
)

type ProtocolHandler interface {
	IsSupported(service Service) bool
	IsPasswordNeeded(service Service) bool
	OpenConnection(service Service)
	Run(service Service, cmd string) (string, string)
	CloseConnection(service Service)
}

type SSHProtocolHandler struct {
	client *ssh.Client
}

func (r *SSHProtocolHandler) IsSupported(service Service) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", service.host, service.port))
	defer conn.Close()
	return err == nil
}

func (r *SSHProtocolHandler) IsPasswordNeeded(service Service) bool {
	return true
}

func (r *SSHProtocolHandler) OpenConnection(service Service) {

	config := &ssh.ClientConfig{
		User: service.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(service.password),
		},
	}

	log.Debug("opening connection to %s:%s: ", service.host, service.port)

	conn, error := ssh.Dial("tcp", fmt.Sprintf("%s:%s", service.host, service.port), config)

	if error == nil {
		r.client = conn
	}
}

func (r *SSHProtocolHandler) Run(service Service, cmd string) (string, string) {

	log.Debug("sending cmd: %s", strings.Replace(cmd, service.sudo, "******", -1))

	session, _ := r.client.NewSession()
	defer session.Close()

	so, _ := session.StdoutPipe()

	stderr, _ := session.CombinedOutput(cmd)
	stdout, _ := ioutil.ReadAll(so)

	sstderr := strings.TrimSpace(string(stderr))
	sstdout := strings.TrimSpace(string(stdout))

	log.Debug("got stdout response %s", sstdout)
	log.Debug("got stderr response %s", sstderr)

	return sstdout, sstderr
}

func (r *SSHProtocolHandler) CloseConnection(service Service) {
	if r.client != nil {
		log.Debug("closing connection to %s:%s: ", service.host, service.port)
		r.client.Close()
	}
}

type WindowsProtocolHandler struct {
}

func (r *WindowsProtocolHandler) IsSupported(service Service) bool {
	return isFileFound("net") || isFileFound("sc.exe")
}

func (r *WindowsProtocolHandler) IsPasswordNeeded(service Service) bool {
	return false
}

func (r *WindowsProtocolHandler) OpenConnection(service Service) {
}

func (r *WindowsProtocolHandler) Run(service Service, cmd string) (string, string) {

	log.Debug("sending cmd: ", cmd)

	s, _ := exec.Command(cmd).CombinedOutput()

	log.Debug("got response ", s)

	return string(s), string(s)
}

func (r *WindowsProtocolHandler) CloseConnection(service Service) {
}
