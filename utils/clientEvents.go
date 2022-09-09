package utils

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

/* func (client *Client) receive() {
	for {
		client.socket.SetReadDeadline(time.Now().Add(time.Second * 30))

		//read an store size of filename
		var filenameSize int64
		err := binary.Read(client.socket, binary.LittleEndian, &filenameSize)
		check(err)

		//create and fill filename buffer with max size filenamesize
		filename := bytes.NewBuffer(make([]byte, 0, filenameSize))
		//copy and store amount of bytes of filename
		bytesRead, err := io.CopyN(filename, client.socket, filenameSize)
		check(err)
		fmt.Printf("Expected %d bytes for filename, read %d bytes\n", filenameSize, bytesRead)

		str := filename.String()
		fmt.Println(strings.TrimLeft(str, " "))

		//read an store size of file
		var filesize int64
		err = binary.Read(client.socket, binary.LittleEndian, &filesize)
		check(err)
		fmt.Printf("Expecting %d bytes in file\n", filesize)

		//create file
		f, err := os.Create(str)
		check(err)
		bytesWritten, err := io.CopyN(f, client.socket, filesize)
		check(err)
		fmt.Printf("Transfer complete, expected %d bytes, wrote %d bytes\n", filesize, bytesWritten)
		if filesize != bytesWritten {
			fmt.Printf("ERROR! File doesn't match expected size!"<\n)
		}
	}
} */

func StartReceiveMode(c string) {
	ch, _ := strconv.Atoi(c)

	connection := getClient(c)
	defer connection.Close()

	// Initialize client with connection and channel, waiting for data
	client := &Client{socket: connection, channel: ch, data: make(chan string)}
	// Send channel to server
	sendChannel(client.socket, c+"receive")

	/* go client.receive()

	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	} */

	// Set deadline for reading
	client.socket.SetReadDeadline(time.Now().Add(time.Second * 30))
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

func StartSendMode(c, path string) {
	ch, err := strconv.Atoi(c)
	check(err)

	connection := getClient(c)
	defer connection.Close()
	// Initialize client with connection and channel, waiting for data
	client := &Client{socket: connection, channel: ch, data: make(chan string)}

	// Send channel to server
	sendChannel(client.socket, c+"send")

	/* for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		check(err)
	} */

	// Get file name
	fileInfo, err := os.Stat(strings.TrimSpace(path))
	check(err)
	fileName := fileInfo.Name()

	sendFileName(client.socket, fileName)
	sendFile(client.socket, path)
}