package utils

import (
	"flag"
	"fmt"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = ":8080"
	SERVER_TYPE = "tcp"
)

func SetFlag(length int, option []string) *int {
	var channel int
	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.IntVar(&channel, "channel", 1, "Channel selected")
	mySet.Parse(option[length-2:])
	flag.Parse()
	return &channel
}

func getClient(c string) net.Conn {
	fmt.Println("Starting client on channel", c)
	// Initiate conn request actively
	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+SERVER_PORT)
	check(err)

	return conn
}

func getServer() net.Listener {
	fmt.Printf("%v Starting server...\n", Now())
	listener, err := net.Listen(SERVER_TYPE, SERVER_HOST+SERVER_PORT)
	check(err)
	return listener
}
