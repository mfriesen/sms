package main

import (
//"code.google.com/p/go.crypto/ssh"
//"fmt"
//"io/ioutil"
//"strings"
)

type WindowsProtocolHandler struct {
}

func (r *WindowsProtocolHandler) OpenConnection(service Service) {
}

func (r *WindowsProtocolHandler) Run(cmd string) string {

	log.Debug("sending cmd: ", cmd)

	s := "" //strings.TrimSpace(string(result))
	log.Debug("got response ", s)

	return s
}

func (r *WindowsProtocolHandler) CloseConnection() {
}
