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
			// Stores new connected client
			manager.clients[connection] = true
			fmt.Printf(NewConn, Now(), connection.channel)
		case connection := <-manager.unregister:
			// Removes disconnected client
			if _, ok := manager.clients[connection]; !ok {
				fmt.Printf(EndConn, Now())
				delete(manager.clients, connection)
			}
		case message := <-manager.file:
			// Stores content of file in clients
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

// Receive client file to send to respective channels
func (manager *ClientManager) manageFiles(client *Client) {
	for {
		// Sets a deadline for Read calls
		client.socket.SetReadDeadline(time.Now().Add(time.Second * 30))

		// Get fileName
		filename, filenameSize, bytesRead := getNameorChannel(client.socket)
		if filename.Len() > 0 {
			fmt.Printf(ExpectedFilename, filenameSize, bytesRead)

			// Get file content
			file := getFile(client.socket)

			// Send file components (Name and content), plus a token to split them
			manager.file <- filename.String() + GoData + file.String()
			// Set destiny channel
			manager.destiny = client.channel
			// Unregister and removeclient
			manager.unregister <- client
		}

	}

}

// Send file to respective clients
func (manager *ClientManager) sendtoClient(client *Client) {
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			} else {
				// Verify that the file is sent to the correct channels
				if client.channel == manager.destiny {
					defer client.socket.Close()
					// USe token to split data into fileName and content
					file := strings.Split(message, GoData)
					fileName, fileData := file[0], file[1]
					// Send filename to client
					sendFileName(client.socket, fileName)

					filesize := int64(len(fileData))
					// Send file size to client
					err := binary.Write(client.socket, binary.LittleEndian, filesize)
					check(err)

					// Send file to client
					bytesWritten, _ := io.WriteString(client.socket, fileData)
					check(err)
					// Validates if file content size is same as bytes written
					compare(int64(bytesWritten), int64(len(fileData)), filesize)

					manager.unregister <- client
				} else {
					// If there isn't channel designed to send
					fmt.Println("There's no channel selected to send file")
				}
			}
		}
	}
}

func StartServerMode() {

	listener := getServer()
	defer listener.Close()

	// Creates manager
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		file:       make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go manager.start()

	for {
		listener, err := listener.Accept()
		check(err)

		// Get channel number
		channel, _, _ := getNameorChannel(listener)
		ch := channel.String()
		op := ch[1:]
		intch, _ := strconv.Atoi(ch[:1])
		check(err)

		// Creates client
		client := &Client{socket: listener, channel: intch, data: make(chan string)}

		// Register client
		manager.register <- client
		if op == "send" {
			go manager.manageFiles(client)
		} else {
			go manager.sendtoClient(client)
		}
	}
}
