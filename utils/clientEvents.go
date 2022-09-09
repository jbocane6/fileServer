package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func (client *Client) receive() {
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
		fmt.Printf("Transfer complete, expected %d bytes, wrote %d bytes", filesize, bytesWritten)
		if filesize != bytesWritten {
			fmt.Printf("ERROR! File doesn't match expected size!")
		}
	}
}

func StartReceiveMode(c string) {
	ch, _ := strconv.Atoi(c)

	connection := getClient(c)
	defer connection.Close()

	client := &Client{socket: connection, channel: ch, data: make(chan string)}
	//connection.Write([]byte(c))

	go client.receive()

	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

func StartSendMode(c, path string) {
	ch, _ := strconv.Atoi(c)

	connection := getClient(c)
	defer connection.Close()

	client := &Client{socket: connection, channel: ch, data: make(chan string)}

	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		check(err)

		fileInfo, err := os.Stat(strings.TrimSpace(message))
		check(err)
		fileName := fileInfo.Name()

		length := int64(len(fileName))
		err = binary.Write(client.socket, binary.LittleEndian, length)
		check(err)

		bytes, err := io.WriteString(client.socket, fileName)
		check(err)
		if bytes != len(fileName) {
			fmt.Printf("Error! Wrote %d bytes but length of name is %d!\n", bytes, length)
		}

		f, err := os.Open(strings.TrimSpace(message))
		check(err)

		stat, err := f.Stat()
		check(err)

		filesize := stat.Size()
		err = binary.Write(client.socket, binary.LittleEndian, filesize)
		check(err)

		bytesWritten, err := io.CopyN(client.socket, f, filesize)
		check(err)
		if bytesWritten != filesize {
			fmt.Printf("Error! Wrote %d bytes but length of file is %d!\n", bytes, stat.Size())
		}
	}
}
