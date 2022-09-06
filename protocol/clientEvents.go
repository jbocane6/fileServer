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
		fl := []byte{}
		if length > 0 {
			os.WriteFile("fakego.png", readMultipleBytes(client.socket, message, fl), 0333)
		}
	}

}

func StartReceiveMode(c string) {
	ch, _ := strconv.Atoi(c)

	connection := getClient(c)
	client := &Client{socket: connection, channel: ch, data: make(chan []byte)}
	connection.Write([]byte(c))

	go client.receive()

	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

func StartSendMode(c, path string) {
	message := make([]byte, 4096)

	connection := getClient(c)
	connection.Write([]byte(c + fileName(path)))
	//connection.Write([]byte(c[4] + "s"))

	val, _ := connection.Read(message)
	fmt.Println(string(message[:val]))

	// start sending file
	wholeFile, _ := os.ReadFile(path)
	connection.Write(wholeFile)
	// start sending file
}

func fileName(path string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("OS. Stat() function execution error, error is:% v \n", err)
		os.Exit(1)
	}
	return fileInfo.Name()
}
