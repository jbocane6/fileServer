package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func sendChannel(socket net.Conn, channel string) {
	// Writes channel length into socket
	err := binary.Write(socket, binary.LittleEndian, int64(len(channel)))
	check(err)

	// Writes channel value into socket
	_, err = io.WriteString(socket, channel)
	check(err)
}

// Get the file name or channel and return value, size and bytes read
func getNameorChannel(socket net.Conn) (*bytes.Buffer, int64, int64) {

	// Get the size of the file name or the channel name
	var size int64
	binary.Read(socket, binary.LittleEndian, &size)

	// Get the content of the file name or the channel name
	data := bytes.NewBuffer(make([]byte, 0, size))
	bytesRead, err := io.CopyN(data, socket, size)
	check(err)

	return data, size, bytesRead
}

func sendFileName(socket net.Conn, fileName string) {

	// Get and send the file name size
	length := int64(len(fileName))
	err := binary.Write(socket, binary.LittleEndian, length)
	check(err)

	// Send the file name content
	bytes, err := io.WriteString(socket, fileName)
	check(err)
	compare(int64(bytes), int64(len(fileName)), length)
	fmt.Printf(ExpectedWriteName, length, bytes)
}
