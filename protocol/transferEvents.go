package protocol

import (
	"fmt"
	"io"
	"net"
	"os"
)

func OpenFl(filePath string) *os.File {
	//Read only open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("OS. Open() function execution error, error is:% v \n", err)
		os.Exit(1)
	}

	return file
}

func SendFile(conn net.Conn, file *os.File) {

	defer file.Close()

	buf := make([]byte, 4096)
	for {
		//Read the data from the local file and write it to the network receiver. How much to read, how much to write
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("sending file completed \n")
			} else {
				fmt.Printf("file. Read() method execution error, error is:% v \n", err)
			}
			return
		}
		//Write to network socket
		conn.Write(buf[:n])
	}
}

func TransferFile(fileName string, connection net.Conn) {
	buf := make([]byte, 4096)
	_, err := connection.Read(buf)
	if err == nil {
		fmt.Printf("Sending to respectively client\n")
	}
	defer connection.Close()
}

func ReceiveFile(conn net.Conn, fileName string) {
	//Create a new file by file name
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("OS. Create() function execution error, error is:% v \n", err)
		return
	}
	defer file.Close()

	//Read data from network and write to local file
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)

		//Write to local file, read and write
		file.Write(buf[:n])
		if err != nil {
			if err == io.EOF {
				fmt.Printf("receive file complete. \n ")
			} else {
				fmt.Printf("conn.read() method execution error, error is:% v \n", err)
			}
			return
		}
	}
}
