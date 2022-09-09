package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
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

		channel, _, _ := getNameorChannel(client.socket)
		ch, _ := strconv.Atoi(channel.String())
		client.channel = ch

		filename, filenameSize, bytesRead := getNameorChannel(client.socket)
		fmt.Printf("Expected %d bytes for filename, read %d bytes\n", filenameSize, bytesRead)

		file := getFile(client.socket)
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

			file := strings.Split(message, "/godata/")
			fileName, fileData := file[0], file[1]
			sendFileName(client.socket, fileName)

			filesize := int64(len(fileData))
			err := binary.Write(client.socket, binary.LittleEndian, filesize)
			check(err)

			bytesWritten, err := io.WriteString(client.socket, fileData)
			check(err)
			compare(int64(bytesWritten), int64(len(fileData)), filesize)
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

		client := &Client{socket: connection, channel: 0, data: make(chan string)}

		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}
