package main

import (
	"fileServer/protocol"
	"flag"
	"fmt"
	"os"
)

func main() {

	//Get command line parameters
	option := os.Args
	//Absolute path to extract file
	if option[1] == "send" {
		if len(option) != 5 {
			fmt.Printf("format: ./client send [path/filename] -channel [channel number]\n")
			return
		}
		var channel int
		mySet := flag.NewFlagSet("", flag.ExitOnError)
		mySet.IntVar(&channel, "channel", 1, "Channel selected")
		mySet.Parse(option[3:])
		flag.Parse()
	}

	connection := protocol.GetClient()
	defer connection.Close()

	path := option[2]
	fileInfo := protocol.CreatePath(path)
	fileName := protocol.SendName(connection, fileInfo)

	if fileName != "" {
		protocol.SendFile(connection, protocol.OpenFl(path))
	}
	//protocol.ReceiveFile(connection, fileName)

}
