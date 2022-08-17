package protocol

import (
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8000"
	SERVER_TYPE = "tcp"
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
	//Initiate connection request actively
	connection, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Printf("net. Dial() function execution error, error is:% v \n", err)
		os.Exit(1)
	}

	return connection
}
