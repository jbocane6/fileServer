package protocol

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func GetServer() net.Listener {
	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")
	return server
}

func GetClient() net.Conn {
	//Initiate conn request actively
	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Printf("net. Dial() function execution error, error is:% v \n", err)
		os.Exit(1)
	}

	return conn
}

func SetFlag(length int, option []string) *int {
	var channel int
	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.IntVar(&channel, "channel", 1, "Channel selected")
	mySet.Parse(option[length-2:])
	flag.Parse()
	return &channel
}
