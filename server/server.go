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
		conn, err := server.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		} else {
			fmt.Println("client connected")
			fileName := protocol.WriteName(conn)
			protocol.ReceiveFile(conn, fileName)
			//protocol.TransferFile(fileName, conn)
		}
	}
}
