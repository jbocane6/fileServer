package utils

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Subscribe(c string) {
	for {
		ch, _ := strconv.Atoi(c)
		connection := getClient(c)

		// Initialize client with connection and channel, waiting for data
		client := &Client{socket: connection, channel: ch, data: make(chan string)}
		// Send channel number to server
		sendChannel(client.socket, c+"receive")

		// Get filename, first size of bytes and then name
		fileName, fileNameSize, bytesRead := getNameorChannel(client.socket)
		// If get a fileName, continue receiving
		if fileName.Len() != 0 {
			// Print filename copy status
			fmt.Printf(ExpectedFilename, fileNameSize, bytesRead)
			str := fileName.String()

			// Get file size
			var fileSize int64
			err := binary.Read(client.socket, binary.LittleEndian, &fileSize)
			check(err)
			// Print file copy status
			fmt.Printf(ExpectingFile, fileSize)

			// Create file using connection, file name and file size
			createFile(client.socket, str, fileSize)
		}
	}
}

func StartReceiveMode(c string) {
	ch, _ := strconv.Atoi(c)
	connection := getClient(c)

	// Initialize client with connection and channel, waiting for data
	client := &Client{socket: connection, channel: ch, data: make(chan string)}
	// Send channel number to server
	sendChannel(client.socket, c+"receive")

	// Get filename, first size of bytes and then name
	fileName, fileNameSize, bytesRead := getNameorChannel(client.socket)
	if fileNameSize != 0 {
		// Print filename copy status
		fmt.Printf(ExpectedFilename, fileNameSize, bytesRead)

		str := fileName.String()

		// Get file size
		var fileSize int64
		err := binary.Read(client.socket, binary.LittleEndian, &fileSize)
		check(err)
		// Print file copy status
		fmt.Printf(ExpectingFile, fileSize)

		// Create file using connection, file name and file size
		createFile(client.socket, str, fileSize)
	}
}

func StartSendMode(c, path string) {
	ch, err := strconv.Atoi(c)
	check(err)

	connection := getClient(c)
	defer connection.Close()
	// Initialize client with connection and channel, waiting for data
	client := &Client{socket: connection, channel: ch, data: make(chan string)}
	// Send channel number to server
	sendChannel(client.socket, c+"send")

	// Get file name
	fileInfo, err := os.Stat(strings.TrimSpace(path))
	check(err)
	fileName := fileInfo.Name()

	sendFileName(client.socket, fileName)
	sendFile(client.socket, path)
	// Show end of connection message
	fmt.Printf(EndConn, Now())
}

func StartClientMode() {
	for {
		fmt.Println("Waiting for option")
		// Read and manage read instructions
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		clientType := strings.Split(strings.TrimSpace(message), " ")
		switch clientType[0] {
		case "subscribe":
			if len(clientType) != 2 {
				fmt.Println(FormatSubscribe)
			} else {
				// Call subscribe method
				Subscribe(clientType[1])
			}
		case "send":
			if len(clientType) != 3 {
				fmt.Println(FormatPublish)
			} else {
				// Call send method
				StartSendMode(clientType[1], clientType[2])
			}
		}
	}

}
