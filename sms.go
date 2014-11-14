package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
	"github.com/howeyc/gopass"
	"github.com/jcelliott/lumber"

	"os"
	"os/exec"
	"os/user"
	"reflect"
	"regexp"
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

	if hasKey(options, "--user") {
		service.user = options["--user"].(string)
	} else {
		usr, _ := user.Current()
		service.user = usr.Username
	}

	if hasKey(options, "<host>") {

		host := options["<host>"].(string)

		rp := regexp.MustCompile("[@:]")
		strs := rp.Split(host, -1)

		if len(strs) == 3 {
			service.host = strs[1]
			service.user = strs[0]
			service.port = strs[2]
		} else if len(strs) == 2 {
			service.host = strs[1]
			service.user = strs[0]
			service.port = DEFAULT_PORT
		} else {
			service.host = host
			service.port = DEFAULT_PORT
		}
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

	if options["restart"] == true {
		service.action = "restart"
	}

	if hasKey(options, "<servicename>") {
		service.name = options["<servicename>"].(string)
	}

	if hasKey(options, "--password") {
		service.password = options["--password"].(string)
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
  sms [options] [user@]<host>[:port] <servicename> restart
  sms [options] [user@]<host>[:port] <servicename> start
  sms [options] [user@]<host>[:port] <servicename> status
  sms [options] [user@]<host>[:port] <servicename> stop

 Options:
  --user=userid  userid
  --password=password  password
  --sudo=sudopw  sudo password
  -h, --help     show help
  -v, --verbose  show debug info
`

	arguments, err := docopt.Parse(usage, argv, true, "0.1", false, exit)

	service = updateOptions(service, arguments)

	return service, err
}

func run(service Service) error {

	var err error
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

			err = protocol.OpenConnection(service)

			if err == nil {
				for _, handler := range handlers {

					handler_supported := handler.IsSupported(protocol)
					log.Debug("checking handler support for %s ... %t", reflect.TypeOf(handler), handler_supported)

					if handler_supported {

						status := ServiceStatusUnknown

						if service.action == "status" {
							status, err = handler.Status(service, protocol)
						} else if service.action == "start" {
							status, err = handler.Start(service, protocol)
						} else if service.action == "stop" {
							status, err = handler.Stop(service, protocol)
						} else if service.action == "restart" {

							status, err = handler.Stop(service, protocol)

							if err == nil {
								status, err = handler.Start(service, protocol)
							}
						}

						if err != nil {
							fmt.Println(fmt.Sprintf("an error ocurred %s", err.Error()))
						} else {
							fmt.Println(fmt.Sprintf("service %s is %s", service.name, ServiceStatus[status]))
						}

						completed = true
						break
					}

				}

				protocol.CloseConnection(service)
			} else {

				completed = true
			}
		}

		if completed {
			break
		}
	}

	return err
}

func main() {

	service, err := usage(os.Args[1:], true)

	if err == nil {
		err = run(service)

		if err != nil {
			fmt.Println(err.Error())
		}
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
