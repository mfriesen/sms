package main

import (
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"io/ioutil"
	"strings"
)

type ProtocolHandler interface {
	OpenConnection(service Service)
	Run(cmd string) string
	CloseConnection()
}

type SSHProtocolHandler struct {
	client *ssh.Client
}

func (r *SSHProtocolHandler) OpenConnection(service Service) {

	config := &ssh.ClientConfig{
		User: service.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(service.password),
		},
	}

	conn, error := ssh.Dial("tcp", fmt.Sprintf("%s:%s", service.host, service.port), config)

	if error != nil {
		panic(error)
	}

	r.client = conn
}

func (r *SSHProtocolHandler) Run(cmd string) string {

	log.Debug("sending cmd: ", cmd)

	session, _ := r.client.NewSession()
	defer session.Close()

	so, error := session.StdoutPipe()
	if error != nil {
		panic(error)
	}

	result, error := session.CombinedOutput(cmd)

	result, error = ioutil.ReadAll(so)
	if error != nil {
		panic(error)
	}

	s := strings.TrimSpace(string(result))
	log.Debug("got response ", s)

	return s
}

func (r *SSHProtocolHandler) CloseConnection() {
	if r.client != nil {
		r.client.Close()
	}
}
