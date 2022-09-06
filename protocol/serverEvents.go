package protocol

import (
	"fmt"
	"strconv"
)

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Printf("%v Added new connection!, channel: %v\n", Now(), connection.channel)
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				fmt.Println("A connection has terminated!")
				delete(manager.clients, connection)
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
		case message := <-manager.fileName:
			fmt.Println("File name to share:", string(message))
		}
	}
}

func (manager *ClientManager) manageFile(client *Client) {
	for {
		message := make([]byte, 4096)
		fl := []byte{}
		manager.file <- readMultipleBytes(client.socket, message, fl)
		manager.destiny <- client.channel
	}
}

func (manager *ClientManager) SendFile(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			} else if client.channel == <-manager.destiny {
				client.socket.Write(message)
			}
			fmt.Printf("%v SENDING: file to channel: %d\n", Now(), client.channel)
		}
	}
}

func StartServerMode() {
	listener := getServer()
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		fileName:   make(chan []byte),
		file:       make(chan []byte),
		destiny:    make(chan int),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
	for {
		connection, error := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		message := make([]byte, 4096)
		c, _ := connection.Read(message)
		if c > 1 {
			name := message[1:]
			manager.fileName <- name
		}
		ch, _ := strconv.Atoi(string(message[:1]))
		connection.Write([]byte("Server accepted connection"))
		client := &Client{socket: connection, channel: ch, data: make(chan []byte)}
		manager.register <- client
		if c == 1 {
			go manager.SendFile(client)
		} else {
			go manager.manageFile(client)
		}
	}
}
