package protocol

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (client *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
		}
	}
}

func StartReceiveMode(c string) {
	ch, _ := strconv.Atoi(c)
	fmt.Println("Starting client on channel ", c)
	connection := GetClient()
	client := &Client{socket: connection, channel: ch}
	connection.Write([]byte(c))
	go client.receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

func StartSendMode(c []string) {
	fmt.Println("Starting client on channel ", c[4])
	message := make([]byte, 4096)
	connection := GetClient()
	connection.Write([]byte(c[4] + "s"))
	val, _ := connection.Read(message)
	fmt.Println(string(message[:val]))
	connection.Write([]byte(c[2]))
	fmt.Println("Sending message: ", c[2])
}
