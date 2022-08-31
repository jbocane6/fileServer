package main

import (
	p "fileServer/protocol"
	"fmt"
	"os"
)

func main() {
	// Get inserted arguments.
	option := os.Args
	// Validate options.
	switch option[1] {
	case "send":
		if len(option) != 5 {
			fmt.Printf("Format: ./client send [path/filename] -channel [channel number]\n")
			return
		}
		p.StartSendMode(option)
	case "receive":
		if len(option) != 4 {
			fmt.Printf("Format: ./client receive -channel [channel number]\n")
			return
		}
		p.StartReceiveMode(option[3])
	}

}
