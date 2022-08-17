package main

import (
	"fileServer/protocol"
	"fmt"
	"os"
)

func main() {
	//activeClients := make(chan protocol.ActiveClient)
	//go protocol.GenerateResponses(activeClients)

	server := protocol.GetServer()
	defer server.Close()
	for {
		connection, err := server.Accept()
		defer connection.Close()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		} else {
			fmt.Println("client connected")
			fileName := protocol.WriteName(connection)
			protocol.ReceiveFile(connection, fileName)
			//protocol.TransferFile(fileName, connection)
		}
	}
}
