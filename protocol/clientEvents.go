package protocol

import (
	"bufio"
	"fmt"
	"io"
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
			for {
				_, err := client.socket.Read(message)
				if err != nil {
					if err == io.EOF {
					} else {
						client.socket.Close()
						break
					}
					break
				}
				fl = append(fl, message...)
			}
			os.WriteFile("fakego.png", fl, 0333)
		}
	}

}

func StartReceiveMode(c string) {
	ch, _ := strconv.Atoi(c)

	connection := GetClient(c)
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
	message := make([]byte, 4096)

	connection := GetClient(c[4])
	/* fileInfo, err := os.Stat(c[2])
	if err != nil {
		fmt.Printf("OS. Stat() function execution error, error is:% v \n", err)
		os.Exit(1)
	}
	fileName := fileInfo.Name()
	connection.Write([]byte(c[4] + fileName)) */
	connection.Write([]byte(c[4] + "s"))

	val, _ := connection.Read(message)
	fmt.Println(string(message[:val]))

	// start sending file
	wholeFile, _ := os.ReadFile(c[2])
	connection.Write(wholeFile)
	// start sending file
}
