package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/howeyc/gopass"
	"github.com/jcelliott/lumber"
	//	"io"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

const DEFAULT_PORT string = "22"

type Service struct {
	user     string
	password string
	host     string
	port     string
	name     string
	action   string
	sudo     string
}

var (
	log = lumber.NewConsoleLogger(lumber.WARN)
)

const (
	ServiceStatusUnknown = iota
	ServiceStatusStopped = iota
	ServiceStatusStarted = iota
)

var ServiceStatus = [...]string{
	"unknown",
	"stopped",
	"started",
}

func updateOptions(service Service, options map[string]interface{}) Service {

	if options["--verbose"] == true {
		log.Level(lumber.DEBUG)
	}

	if options["start"] == true {
		service.action = "start"
	}

	if options["status"] == true {
		service.action = "status"
	}

	if options["stop"] == true {
		service.action = "stop"
	}

	if hasKey(options, "<servicename>") {
		service.name = options["<servicename>"].(string)
	}

	if hasKey(options, "--sudo") {
		service.sudo = options["--sudo"].(string)

		if service.sudo == "" {
			fmt.Printf(fmt.Sprintf("[sudo] password for %s: ", service.user))
			pass := gopass.GetPasswd()
			service.sudo = string(pass)
		}
	}

	return service
}

func hasKey(m map[string]interface{}, key string) bool {

	var exists bool

	if _, ok := m[key]; ok && m[key] != nil {
		exists = true
	} else {
		exists = false
	}

	return exists
}

func usage(argv []string, exit bool) (Service, error) {

	var service Service
	var err error

	usage := `Service Monitoring System
Usage:
  sms [options] <user>@<host> <servicename> start
  sms [options] <user>@<host> <servicename> status
  sms [options] <user>@<host> <servicename> stop

 Options:
  --sudo=sudopw  sudo password
  -h, --help     show help
  -v, --verbose  show debug info
`

	arguments, err := docopt.Parse(usage, argv, true, "0.1", false, exit)

	if _, ok := arguments["<user>@<host>"]; err == nil && ok {

		var userhost string = arguments["<user>@<host>"].(string)

		if strings.Contains(userhost, "@") {

			rp := regexp.MustCompile("[@:]")
			split := rp.Split(userhost, -1)

			service, err = userHostUsage(split, exit)
		} else {
			err = errors.New("missing '@'")
		}

		if err != nil { // invalid <user>@<host>
			docopt.Parse(usage, []string{}, true, "", false, exit)
		}
	}

	service = updateOptions(service, arguments)

	return service, err
}

func userHostUsage(argv []string, exit bool) (Service, error) {
	var service Service

	usage := `Usage:
	  sm <user> <host> [<port>]`

	arguments, error := docopt.Parse(usage, argv, true, "", false, exit)

	if error == nil {

		if arguments["<port>"] == nil {
			arguments["<port>"] = DEFAULT_PORT
		}

		service = Service{user: arguments["<user>"].(string),
			host: arguments["<host>"].(string),
			port: arguments["<port>"].(string)}
	}

	return service, error
}

func run(service Service) {

	completed := false
	protocols := [...]ProtocolHandler{
		ProtocolHandler(&SSHProtocolHandler{}),
		ProtocolHandler(&WindowsProtocolHandler{}),
	}

	handlers := [...]ServiceHandler{
		ServiceHandler(&ServiceExecServiceHandler{}),
		ServiceHandler(&ScExecServiceHandler{}),
		ServiceHandler(&SambaServiceHandler{}),
	}

	for _, protocol := range protocols {

		supported := protocol.IsSupported(service)
		log.Debug("checking protocol support for %s ... is supported %t", reflect.TypeOf(protocol), supported)

		if supported {

			if protocol.IsPasswordNeeded(service) && service.password == "" {

				fmt.Printf(fmt.Sprintf("%s@%s's Password: ", service.user, service.host))
				pass := gopass.GetPasswd()
				service.password = string(pass)
			}

			protocol.OpenConnection(service)

			for _, handler := range handlers {

				handler_supported := handler.IsSupported(protocol)
				log.Debug("checking handler support for %s ... %t", reflect.TypeOf(handler), handler_supported)

				if handler_supported {

					status := ServiceStatusUnknown

					if service.action == "status" {
						status = handler.Status(service, protocol)
					} else if service.action == "start" {
						status = handler.Start(service, protocol)
					} else if service.action == "stop" {
						status = handler.Stop(service, protocol)
					}

					fmt.Println(fmt.Sprintf("service %s is %s", service.name, ServiceStatus[status]))

					completed = true
					break
				}

			}

			protocol.CloseConnection(service)
		}

		if completed {
			break
		}
	}
}

func main() {

	service, err := usage(os.Args[1:], true)

	if err == nil {
		run(service)
	}
}

func isFileFound(file string) bool {

	_, error := exec.LookPath(file)
	supported := error == nil

	if error != nil {
		log.Debug("cannot find '%s' executable ", file)
	}

	return supported
}
