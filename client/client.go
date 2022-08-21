package main

import (
	"fileServer/protocol"
	"fmt"
	"os"
)

func main() {

	// Get inserted arguments.
	option := os.Args
	// Validate option.
	switch option[1] {
	case "send":
		if len(option) != 5 {
			fmt.Printf("Format: ./client send [path/filename] -channel [channel number]\n")
			return
		}
		fmt.Printf("Channel choose: %d\n", *protocol.SetFlag(len(option), option))
		conn := protocol.GetClient()
		defer conn.Close()
		protocol.SendFile(conn, option[2])
	case "receive":
		if len(option) != 4 {
			fmt.Printf("Format: ./client receive -channel [channel number]\n")
			return
		}
		fmt.Printf("Channel choose: %d\n", *protocol.SetFlag(len(option), option))
		conn := protocol.GetClient()
		defer conn.Close()
	}
}
