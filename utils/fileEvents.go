package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func createFile(socket net.Conn, str string, fileSize int64) {
	// Creates file and stores in f
	f, err := os.Create(str)
	check(err)

	// Copy content of file from connection
	bytesWritten, err := io.CopyN(f, socket, fileSize)
	check(err)
	fmt.Printf(MsgTransfer, Now(), fileSize, bytesWritten)
	if fileSize != bytesWritten {
		fmt.Println(ErrFileSize)
	}
}

func sendFile(socket net.Conn, path string) {
	// Get file content
	f, err := os.Open(strings.TrimSpace(path))
	check(err)

	// Get the file structure
	stat, err := f.Stat()
	check(err)

	// Get length of file
	filesize := stat.Size()
	// Send file length to server
	err = binary.Write(socket, binary.LittleEndian, filesize)
	check(err)

	// Send file content to server
	bytesWritten, err := io.CopyN(socket, f, filesize)
	check(err)
	compare(bytesWritten, filesize, stat.Size())
}

func getFile(socket net.Conn) *bytes.Buffer {
	// Read an store size of file
	var filesize int64
	err := binary.Read(socket, binary.LittleEndian, &filesize)
	check(err)
	//fmt.Printf("Expecting %d bytes in file\n", filesize)

	// Read and return file content
	file := bytes.NewBuffer(make([]byte, 0, filesize))
	bytesFile, err := io.CopyN(file, socket, filesize)
	check(err)
	fmt.Printf(ExpectedFile, filesize, bytesFile)
	return file
}
