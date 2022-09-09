package main

import (
	u "fileServer/utils"
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
		ch := *u.SetFlag(len(option), option)
		//u.StartClientMode(strconv.Itoa(ch))
		u.StartSendMode(strconv.Itoa(ch), option[2])
	case "receive":
		if len(option) != 4 {
			fmt.Printf("Format: ./client receive -channel [channel number]\n")
			return
		}
		ch := *u.SetFlag(len(option), option)
		//u.StartClientMode(strconv.Itoa(ch))
		u.StartReceiveMode(strconv.Itoa(ch))
	}
}
