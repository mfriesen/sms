package main

import (
	//"fmt"
	"errors"
	"github.com/docopt/docopt-go"
	"os"
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

const (
	ServiceStatusUnknown = iota
	ServiceStatusStopped = iota
	ServiceStatusStarted = iota
)

func usage(argv []string, exit bool) (Service, error) {

	var service Service
	var err error

	usage := `Usage:
	  sm [options] <user>@<host>

	  options:
		-h, --help`

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

func main() {

	_, err := usage(os.Args[1:], true)

	if err == nil {

	}
}
