package utils

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = ":8080"
	SERVER_TYPE = "tcp"
)

func SetFlag(length int, option []string) *int {
	var channel int
	// Returns a new, empty flag set with a specified name and error handling property
	mySet := flag.NewFlagSet("", flag.ExitOnError)
	// Creates an int flag with name, default value, and usage string
	mySet.IntVar(&channel, "channel", 1, "Channel selected")
	// Parses flag definitions from the argument list
	mySet.Parse(option[length-2:])
	flag.Parse()
	return &channel
}

// Set new Client address
func getClient(c string) net.Conn {
	if strings.TrimSpace(c) != "" {
		fmt.Printf(StartClient, Now(), c)
	} else {
		fmt.Printf(StartClientnoCh, Now())
	}
	// Initiate conn request actively
	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+SERVER_PORT)
	check(err)

	return conn
}

// Set Server address
func getServer() net.Listener {
	fmt.Printf(StartServer, Now())
	// Initiate conn listener
	listener, err := net.Listen(SERVER_TYPE, SERVER_HOST+SERVER_PORT)
	check(err)
	return listener
}
