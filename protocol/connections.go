package protocol

import (
	"flag"
	"fmt"
	"net"
)

func SetFlag(length int, option []string) *int {
	var channel int
	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.IntVar(&channel, "channel", 1, "Channel selected")
	mySet.Parse(option[length-2:])
	flag.Parse()
	return &channel
}

func GetClient() net.Conn {
	// Initiate conn request actively
	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println(err)
	}
	// conn.Write([]byte(strconv.Itoa(channel)))

	return conn
}
