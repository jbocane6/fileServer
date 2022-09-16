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
	if len(option) > 1 {
		switch option[1] {
		case "send":
			if len(option) != 5 {
				fmt.Printf(u.FormatSend)
				return
			}
			// Set channel flag
			ch := *u.SetFlag(len(option), option)
			// Start send file mode
			u.StartSendMode(strconv.Itoa(ch), option[2])
		case "receive":
			if len(option) != 4 {
				fmt.Printf(u.FormatReceive)
				return
			}
			// Set channel flag
			ch := *u.SetFlag(len(option), option)
			// Start receive file mode
			u.StartReceiveMode(strconv.Itoa(ch))
		}
	} else {
		// If no parameters are given
		u.StartClientMode()
	}
}
