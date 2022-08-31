package main

import (
	p "fileServer/protocol"
)

func main() {
	p.StartServerMode()
	/* flagMode := flag.String("mode", "server", "start in client or server mode")
	flag.Parse()
	if strings.ToLower(*flagMode) == "server" {
		startServerMode()
	} else {
		startClientMode()
	} */
}
