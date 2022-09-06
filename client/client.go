package main

import (
	p "fileServer/protocol"
	"fmt"
	"os"
	"strconv"
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
		ch := *p.SetFlag(len(option), option)
		p.StartSendMode(strconv.Itoa(ch), option[2])
	case "receive":
		if len(option) != 4 {
			fmt.Printf("Format: ./client receive -channel [channel number]\n")
			return
		}
		ch := *p.SetFlag(len(option), option)
		p.StartReceiveMode(strconv.Itoa(ch))
	}
}
