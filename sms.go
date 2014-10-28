package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/howeyc/gopass"
	"github.com/jcelliott/lumber"
	"os"
	"regexp"
	//"runtime"
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

	if _, ok := options["<servicename>"]; ok {
		service.name = options["<servicename>"].(string)
	}

	return service
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

	ssh := SSHProtocolHandler{}
	handler := ProtocolHandler(&ssh)

	r := ServiceHandler(&LinuxServiceHandler{handler: handler})

	r.Connect(service)
	r.Status(service)
	r.Disconnect(service)
}

func main() {

	//fmt.Println("OS VERSION ", runtime.GOOS)
	service, err := usage(os.Args[1:], true)

	fmt.Printf("%s@%s's password:", service.user, service.host)

	pass := gopass.GetPasswd()
	service.password = string(pass)

	if err == nil {
		run(service)
	}
}
