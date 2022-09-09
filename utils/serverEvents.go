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
		// Stores new connected client
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Printf(NewConn, Now(), connection.channel)
		// Removes disconnected client
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				fmt.Printf(EndConn, Now())
				delete(manager.clients, connection)
			}
		// Stores content of file in clients with compatible channel
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
	for {
		client.socket.SetReadDeadline(time.Now().Add(time.Second * 30))

		filename, filenameSize, bytesRead := getNameorChannel(client.socket)
		fmt.Printf(ExpectedFilename, filenameSize, bytesRead)

		file := getFile(client.socket)
		//fmt.Println("RECEIVED: " + string(message))
		manager.file <- filename.String() + GoData + file.String()
		manager.destiny <- client.channel
		/* manager.unregister <- client
		client.socket.Close() */
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			} else if client.channel == <-manager.destiny {
				file := strings.Split(message, GoData)
				fileName, fileData := file[0], file[1]
				sendFileName(client.socket, fileName)

				filesize := int64(len(fileData))
				err := binary.Write(client.socket, binary.LittleEndian, filesize)
				check(err)

				bytesWritten, err := io.WriteString(client.socket, fileData)
				check(err)
				compare(int64(bytesWritten), int64(len(fileData)), filesize)

				/* manager.unregister <- client
				client.socket.Close() */
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
		listener, err := listener.Accept()
		check(err)

		channel, _, _ := getNameorChannel(listener)
		ch := channel.String()
		op := ch[1:]
		intch, err := strconv.Atoi(ch[:1])
		check(err)
		fmt.Println(op)

		client := &Client{socket: listener, channel: intch, data: make(chan string)}

		manager.register <- client
		if op == "send" {
			go manager.receive(client)
		} else {
			go manager.send(client)
		}
	}
}
