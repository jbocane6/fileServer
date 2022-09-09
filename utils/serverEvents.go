package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"
)

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has terminated!")
			}
		case message := <-manager.file:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	defer client.socket.Close()
	for {
		client.socket.SetReadDeadline(time.Now().Add(time.Second * 30))

		//read an store size of filename
		var filenameSize int64 //size of fileName
		err := binary.Read(client.socket, binary.LittleEndian, &filenameSize)
		check(err)

		//create and fill filename buffer with max size filenamesize
		filename := bytes.NewBuffer(make([]byte, 0, filenameSize)) // filename will be a slice of bytes that store the fileName
		//copy and store amount of bytes of filename
		bytesRead, err := io.CopyN(filename, client.socket, filenameSize) //bytes fileName
		check(err)
		fmt.Printf("Expected %d bytes for filename, read %d bytes\n", filenameSize, bytesRead)

		/* str := strings.TrimLeft(filename.String(), " ")
		fmt.Println(str) */

		//read an store size of file
		var filesize int64 //size of file
		err = binary.Read(client.socket, binary.LittleEndian, &filesize)
		check(err)
		//fmt.Printf("Expecting %d bytes in file\n", filesize)

		file := bytes.NewBuffer(make([]byte, 0, filesize))
		bytesFile, err := io.CopyN(file, client.socket, filesize)
		check(err)
		fmt.Printf("Expected %d bytes for file, read %d bytes\n", filesize, bytesFile)
		//fmt.Println("RECEIVED: " + string(message))
		manager.file <- filename.String() + "/godata/" + file.String()
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			//client.socket.Write(message)
			fmt.Println(len(message))
			file := strings.Split(message, "/godata/")
			fileName := file[0]
			fileData := file[1]
			length := int64(len(fileName))
			err := binary.Write(client.socket, binary.LittleEndian, length)
			check(err)

			bytes, err := io.WriteString(client.socket, fileName)
			check(err)
			if bytes != len(fileName) {
				fmt.Printf("Error! Wrote %d bytes but length of name is %d!\n", bytes, length)

			}

			filesize := int64(len(file[1]))
			err = binary.Write(client.socket, binary.LittleEndian, filesize)
			check(err)

			bytesWritten, err := io.WriteString(client.socket, fileData)
			check(err)
			if bytesWritten != int(filesize) {
				fmt.Printf("Error! Wrote %d bytes but length of file is %d!\n", bytesWritten, filesize)
			}
		}

	}
}

func StartServerMode() {

	listener := getServer()
	defer listener.Close()

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		file:       make(chan string),
		destiny:    make(chan int),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go manager.start()

	for {
		connection, err := listener.Accept()
		check(err)
		client := &Client{socket: connection, data: make(chan string)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}
