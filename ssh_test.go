package main

import (
//"code.google.com/p/go.crypto/ssh"
//"fmt"
//"code.google.com/p/go.crypto/ssh"
//"net"
//	"github.com/maraino/go-mock"
//"testing"
)

/*
const packageVersion = "SSH-2.0-Go"

func testClientVersion(t *testing.T, config *ssh.ClientConfig, expected string) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	receivedVersion := make(chan string, 1)
	go func() {
		version, err := readVersion(serverConn)
		if err != nil {
			receivedVersion <- ""
		} else {
			receivedVersion <- string(version)
		}
		serverConn.Close()
	}()
	ssh.NewClientConn(clientConn, "", config)
	actual := <-receivedVersion
	if actual != expected {
		t.Fatalf("got %s; want %s", actual, expected)
	}
}

func TestCustomClientVersion(t *testing.T) {
	version := "Test-Client-Version-0.0"
	testClientVersion(t, &ssh.ClientConfig{ClientVersion: version}, version)
}

func TestDefaultClientVersion(t *testing.T) {
	testClientVersion(t, &ssh.ClientConfig{}, packageVersion)
}
*/
/*
type MySshClient struct {
	mock.Mock
}

func (c *MySshClient) Dial(network, addr string, config *ClientConfig) (*Client, error) {
	fmt.Println("DIALING>.... FAKE")
	c := &MySshClient{}
	return c
}

// test StartService 'OK'
func TestStartService01(t *testing.T) {
	// given

	// when
	StartService()

	// then
	//	if err == nil {
	//		t.Error("Expected Errors, got none")
	//	}
}
*/
