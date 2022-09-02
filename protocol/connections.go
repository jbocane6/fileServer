package protocol

import (
	"flag"
	"fmt"
	"net"
)

func SetFlag(length int, option []string) *int {
	var channel int
	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.IntVar(&channel, "channel", 1, "Channel selected")
	mySet.Parse(option[length-2:])
	flag.Parse()
	return &channel
}

func GetClient(c string) net.Conn {
	fmt.Println("Starting client on channel ", c)
	// Initiate conn request actively
	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+SERVER_PORT)
	if err != nil {
		fmt.Println(err)
	}
	// conn.Write([]byte(strconv.Itoa(channel)))

	return conn
}

func GetServer() net.Listener {
	fmt.Printf("%v Starting server...\n", Now())
	listener, error := net.Listen(SERVER_TYPE, SERVER_HOST+SERVER_PORT)
	if error != nil {
		fmt.Println(error)
	}
	return listener
}
