package utils

import (
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
		// Send channel to server
		sendChannel(client.socket, c+"receive")

		// Get filename, first size of bytes and then name
		fileName, fileNameSize, bytesRead := getNameorChannel(client.socket)
		fmt.Printf(ExpectedFilename, fileNameSize, bytesRead)

		str := fileName.String()

		// Get file size
		var fileSize int64
		err := binary.Read(client.socket, binary.LittleEndian, &fileSize)
		check(err)
		fmt.Printf(ExpectingFile, fileSize)

		// Create file using connection, file name and file size
		createFile(client.socket, str, fileSize)
	}
}

func StartReceiveMode(c string) {

	ch, _ := strconv.Atoi(c)

	connection := getClient(c)

	// Initialize client with connection and channel, waiting for data
	client := &Client{socket: connection, channel: ch, data: make(chan string)}
	// Send channel to server
	sendChannel(client.socket, c+"receive")

	// Set deadline for reading
	//client.socket.SetReadDeadline(time.Now().Add(time.Second * 30))
	// Get filename, first size of bytes and then name
	fileName, fileNameSize, bytesRead := getNameorChannel(client.socket)
	if fileNameSize != 0 {
		fmt.Printf(ExpectedFilename, fileNameSize, bytesRead)

		str := fileName.String()

		// Get file size
		var fileSize int64
		err := binary.Read(client.socket, binary.LittleEndian, &fileSize)
		check(err)
		fmt.Printf(ExpectingFile, fileSize)

		// Create file using connection, file name and file size
		createFile(client.socket, str, fileSize)
	}
}

func StartSendMode(c, path string) {
	ch, err := strconv.Atoi(c)
	check(err)

	connection := getClient(c)
	// Initialize client with connection and channel, waiting for data
	client := &Client{socket: connection, channel: ch, data: make(chan string)}
	// Send channel to server
	sendChannel(client.socket, c+"send")

	// Get file name
	fileInfo, err := os.Stat(strings.TrimSpace(path))
	check(err)
	fileName := fileInfo.Name()

	sendFileName(client.socket, fileName)
	sendFile(client.socket, path)

}
