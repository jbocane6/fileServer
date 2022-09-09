package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func sendChannel(socket net.Conn, channel int) {
	ch := strconv.Itoa(channel)
	err := binary.Write(socket, binary.LittleEndian, int64(len(ch)))
	check(err)
	_, err = io.WriteString(socket, ch)
	check(err)
}

func getNameorChannel(socket net.Conn) (*bytes.Buffer, int64, int64) {
	var size int64
	err := binary.Read(socket, binary.LittleEndian, &size)
	check(err)
	data := bytes.NewBuffer(make([]byte, 0, size))
	bytesRead, err := io.CopyN(data, socket, size)
	check(err)
	return data, size, bytesRead
}

func sendFileName(socket net.Conn, fileName string) {

	length := int64(len(fileName))
	err := binary.Write(socket, binary.LittleEndian, length)
	check(err)

	bytes, err := io.WriteString(socket, fileName)
	check(err)
	compare(int64(bytes), int64(len(fileName)), length)
}

func getFile(socket net.Conn) *bytes.Buffer {
	// read an store size of file
	var filesize int64 //size of file
	err := binary.Read(socket, binary.LittleEndian, &filesize)
	check(err)
	//fmt.Printf("Expecting %d bytes in file\n", filesize)

	file := bytes.NewBuffer(make([]byte, 0, filesize))
	bytesFile, err := io.CopyN(file, socket, filesize)
	check(err)
	fmt.Printf("Expected %d bytes for file, read %d bytes\n", filesize, bytesFile)
	return file
}

func createFile(socket net.Conn, str string, fileSize int64) {
	f, err := os.Create(str)
	check(err)
	bytesWritten, err := io.CopyN(f, socket, fileSize)
	check(err)
	fmt.Printf("Transfer complete, expected %d bytes, wrote %d bytes\n", fileSize, bytesWritten)
	if fileSize != bytesWritten {
		fmt.Printf("ERROR! File doesn't match expected size!\n")
	}
}

func sendFile(socket net.Conn, path string) {
	f, err := os.Open(strings.TrimSpace(path))
	check(err)

	stat, err := f.Stat()
	check(err)

	filesize := stat.Size()
	err = binary.Write(socket, binary.LittleEndian, filesize)
	check(err)

	bytesWritten, err := io.CopyN(socket, f, filesize)
	check(err)
	compare(bytesWritten, filesize, stat.Size())
}
