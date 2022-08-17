package protocol

import (
	"fmt"
	"net"
	"os"
)

func CreatePath(path string) os.FileInfo {
	//Get file properties
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("OS. Stat() function execution error, error is:% v \n", err)
		os.Exit(1)
	}
	return fileInfo
}

func SendName(connection net.Conn, fileInfo os.FileInfo) string {
	// Send file name to server
	_, err := connection.Write([]byte(fileInfo.Name()))

	//Read server postback data
	buf := make([]byte, 4096)
	n, err := connection.Read(buf)
	if err != nil {
		fmt.Printf("connection.read (buf) method execution error, error is:% v \n", err)
		os.Exit(1)
	}
	return string(buf[:n])
}

func WriteName(conn net.Conn) string {
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("conn.read() method execution error, error is:% v \n", err)
		os.Exit(1)
	}
	fileName := string(buf[:n])

	//Write back OK to the sender
	conn.Write([]byte(fileName))

	return fileName
}
