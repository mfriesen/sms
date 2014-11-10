package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"code.google.com/p/go.crypto/ssh"
)

type ProtocolHandler interface {
	IsSupported(service Service) bool
	IsPasswordNeeded(service Service) bool
	OpenConnection(service Service)
	Run(service Service, cmd string) (string, error)
	CloseConnection(service Service)
}

type SSHProtocolHandler struct {
	client *ssh.Client
}

func (r *SSHProtocolHandler) IsSupported(service Service) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", service.host, service.port))

	supported := err == nil
	if supported {
		defer conn.Close()
	}

	return supported
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
	// TODO handle wrong username / password
	log.Debug("opening connection to %s:%s: ", service.host, service.port)

	conn, error := ssh.Dial("tcp", fmt.Sprintf("%s:%s", service.host, service.port), config)

	if error == nil {
		r.client = conn
	}
}

func (r *SSHProtocolHandler) Run(service Service, cmd string) (string, error) {

	var stdout, stderr, response bytes.Buffer

	cmdString := cmd
	if service.sudo != "" {
		cmdString = strings.Replace(cmd, service.sudo, "******", -1)
	}

	log.Debug("sending cmd: %s", cmdString)

	session, _ := r.client.NewSession()
	defer session.Close()

	session.Stdout = &stdout
	session.Stderr = &stderr

	in, _ := session.StdinPipe()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: %s", err)
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: %s", err)
	}

	// login text
	bytes, err := r.ReadBuffer(&stdout)
	//response.WriteString(string(bytes))

	fmt.Fprintln(in, cmd)
	bytes, err = r.ReadBuffer(&stdout)
	response.WriteString(string(bytes))

	log.Debug("receive %s", response.String())

	if err == io.EOF {
		err = nil
	}

	if err != nil {
		log.Debug("received error %s", err.Error())
	}

	return response.String(), err
}

func (r *SSHProtocolHandler) ReadBuffer(stdout *bytes.Buffer) ([]byte, error) {

	len := -1
	count := 0
	delay := 100

	for {
		time.Sleep(time.Duration(delay) * time.Millisecond)

		stdoutLen := stdout.Len()

		if stdoutLen == len && count > 0 {
			return stdout.ReadBytes(0)
		} else {
			len = stdoutLen
		}

		delay *= 2
		count += 1

		if count > 10 {
			return nil, errors.New("ReadBuffer timeout")
		}
	}
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
	return strings.Contains(runtime.GOOS, "windows")
}

func (r *WindowsProtocolHandler) IsPasswordNeeded(service Service) bool {
	return false
}

func (r *WindowsProtocolHandler) OpenConnection(service Service) {
}

func (r *WindowsProtocolHandler) Run(service Service, cmd string) (string, error) {

	log.Debug("sending cmd: ", cmd)

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	s, err := exec.Command(head, parts...).CombinedOutput()
	//s, err := exec.Command("sh", "-c", cmd).CombinedOutput()

	log.Debug("got response ", string(s))
	log.Debug("got error ", err)

	return string(s), err
}

func (r *WindowsProtocolHandler) CloseConnection(service Service) {
}
